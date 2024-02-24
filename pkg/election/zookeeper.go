package election

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

const (
	persistentNodePath = "/election"
	ephemeralNodePath  = persistentNodePath + "/candidate"
)

// Zookeeper implements the Elector interface using Apache Zookeeper.
type Zookeeper struct {
	// Servers is the list of addresses of Zookeeper nodes.
	Servers []string

	conn   *zk.Conn
	myPath string
}

func (z *Zookeeper) Init() error {
	// Attempt connection with Zookeeper.
	conn, _, err := zk.Connect(z.Servers, time.Second*5, zk.WithLogger(zooLogger{}))
	if err != nil {
		return fmt.Errorf("error in zk.Connect call: %w", err)
	}

	// Persist the connection.
	z.conn = conn

	// Create the persistent zNode for the election.
	if _, err := z.conn.Create(persistentNodePath, nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
		if !errors.Is(err, zk.ErrNodeExists) {
			return fmt.Errorf("error while creating persistent zNode: %w", err)
		}
		// Persistent zNode already exists.
	}

	// Create the ephemeral-sequential zNode.
	path, err := z.conn.CreateProtectedEphemeralSequential(ephemeralNodePath, nil, zk.WorldACL(zk.PermAll))
	if err != nil {
		return fmt.Errorf("error while creating ephemeral-sequential zNode: %w", err)
	}

	// Persist own path.
	z.myPath = strings.TrimPrefix(path, persistentNodePath+"/")
	fmt.Println("INFO: My path:", z.myPath)
	return nil
}

func (z *Zookeeper) Participate() (<-chan struct{}, <-chan error) {
	// Create the required channels.
	leaderChan, errChan := make(chan struct{}), make(chan error)

	go func() {
		for {
			// Get all the children of the persistent node.
			children, _, err := z.conn.Children(persistentNodePath)
			if err != nil {
				errChan <- fmt.Errorf("failed to get persistent node children: %w", err)
				continue
			}

			// Map the child nodes to their parsed sequence numbers.
			childrenSequences, err := mapESNodeSequence(children)
			if err != nil {
				errChan <- fmt.Errorf("failed to parse children sequence: %w", err)
				continue
			}

			// Sort the children based on their sequence number.
			sort.SliceStable(children, func(i, j int) bool {
				return childrenSequences[children[i]] < childrenSequences[children[j]]
			})

			// Find this position of this node's sequence.
			myPosition := sort.Search(len(children), func(i int) bool {
				return children[i] == z.myPath
			})

			fmt.Printf("INFO: All children: %+v\n", children)
			fmt.Printf("INFO: My position: %d\n", myPosition)

			// If this node is the first child, assume leadership.
			if myPosition == 0 {
				leaderChan <- struct{}{}
				break
			}

			// Get the full path of the node above.
			upperNodePath := children[myPosition-1]
			upperNodeFullPath := persistentNodePath + "/" + upperNodePath

			fmt.Printf("INFO: Awaiting deletion of: %s\n", upperNodePath)

			// Await the deletion of upper node.
			if err := z.awaitDeletion(upperNodeFullPath); err != nil {
				errChan <- fmt.Errorf("error while waiting for node deletion: %w", err)
				continue
			}

			fmt.Printf("INFO: %s deleted\n", upperNodePath)
		}

		// Channel cleaup.
		close(leaderChan)
		close(errChan)
	}()

	return leaderChan, errChan
}

// awaitDeletion blocks until the node at the given path is deleted.
func (z *Zookeeper) awaitDeletion(path string) error {
	// Set a watch on the given node.
	exists, _, emitter, err := z.conn.ExistsW(path)
	if err != nil {
		return fmt.Errorf("error while watching zNode %s: %w", path, err)
	}

	// If node doesn't exist.
	if !exists {
		return fmt.Errorf("zNode %s does not exist", path)
	}

	// Keep listening for events.
	for {
		// If the event type is node-deletion, break inifinite loop.
		if event := <-emitter; event.Type == zk.EventNodeDeleted {
			break
		}
	}

	return nil
}
