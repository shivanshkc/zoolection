package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/go-zookeeper/zk"
)

// createNodes function creates a persistent znode if it does not already exist,
// and creates a sequential ephemeral znode under the persistent one to start this
// service's participation in the election.
func createNodes(conn *zk.Conn) (string, error) {
	// Create the persistent zNode for the election.
	if _, err := conn.Create("/election", nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
		if !errors.Is(err, zk.ErrNodeExists) {
			return "", fmt.Errorf("error while creating persistent zNode: %w", err)
		}
		// Persistent zNode already exists.
	}

	// Create the ephemeral-sequential zNode.
	path, err := conn.Create("/election/candidate", nil, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		return "", fmt.Errorf("error while creating ephemeral-sequential zNode: %w", err)
	}

	return path, nil
}

// awaitVictory function blocks until the sequence number of this service's znode is the smallest
// of all sequence numbers, after which the service can assume leadership of the cluster.
func awaitVictory(conn *zk.Conn, myNodePath string) {
	for {
		// Get all the children of the persistent node.
		children, _, err := conn.Children("/election")
		if err != nil {
			fmt.Println("ERROR: failed to get persistent node children:", err)
			continue
		}

		// Sort the children.
		slices.SortStableFunc(children, func(a, b string) int {
			// Ignoring errors for brevity.
			aSequence, _ := parseSequence(a, "candidate")
			bSequence, _ := parseSequence(b, "candidate")

			if aSequence < bSequence {
				return -1
			}
			if aSequence > bSequence {
				return 1
			}
			return 0
		})

		// Find own position in the sorted children list.
		// For very large systems, this should be replaced with binary search.
		var myPosition int
		for i, elem := range children {
			if myNodePath == "/election/"+elem {
				myPosition = i
				break
			}
		}

		fmt.Println("All children:", children)
		fmt.Println("My node:", myNodePath)
		fmt.Println("My rank:", myPosition)

		// If this node is the first child, assume leadership.
		if myPosition == 0 {
			return
		}

		// Get the full path of the node above.
		upperNodePath := children[myPosition-1]
		upperNodeFullPath := "/election/" + upperNodePath

		fmt.Printf("INFO: Awaiting deletion of: %s\n", upperNodePath)

		// Await the deletion of upper node.
		if err := awaitDeletion(conn, upperNodeFullPath); err != nil {
			fmt.Println("ERROR: error while waiting for node deletion:", err)
			continue
		}

		fmt.Printf("INFO: %s deleted\n", upperNodePath)
	}
}

// awaitDeletion blocks until the znode at the given path is deleted.
func awaitDeletion(conn *zk.Conn, path string) error {
	// Set a watch on the given node.
	exists, _, emitter, err := conn.ExistsW(path)
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
