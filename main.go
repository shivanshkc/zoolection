package main

import (
	"fmt"

	"github.com/shivanshkc/zoolection/pkg/election"
)

func main() {
	zook := &election.Zookeeper{}

	// Initial setup.
	if err := zook.Init("localhost:2181", "localhost:2182", "localhost:2183"); err != nil {
		panic(fmt.Errorf("error in zookeeper.Init call: %w", err))
	}

	go func() {
		// Start participating in the election.
		zook.Participate()
		fmt.Println("Elected as leader.")
	}()

	// Block forever.
	select {}
}
