package main

import (
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
)

func main() {
	servers := []string{"localhost:2181", "localhost:2182", "localhost:2183"}

	// Attempt connection with Zookeeper.
	conn, _, err := zk.Connect(servers, time.Second*5)
	if err != nil {
		panic(err)
	}

	// Create the required nodes for election.
	myNodePath, err := createNodes(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println("My node path:", myNodePath)

	// Block until elected as leader.
	awaitVictory(conn, myNodePath)

	fmt.Println("Elected as leader.")

	// Do leader stuff.
	select {}
}
