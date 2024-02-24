package main

import (
	"fmt"

	"github.com/shivanshkc/zoolection/pkg/election"
)

func main() {
	zookeeper := &election.Zookeeper{
		Servers: []string{
			"localhost:2181",
			"localhost:2182",
			"localhost:2183",
		},
	}

	// Initial setup.
	if err := zookeeper.Init(); err != nil {
		panic(fmt.Errorf("error in zookeeper.Init call: %w", err))
	}

	// Start participating in the election.
	electedChan, errorChan := zookeeper.Participate()

	// Goroutine to report successful election.
	go func() {
		<-electedChan
		fmt.Println("Elected as leader.")
	}()

	// Goroutine to report errors.
	go func() {
		for err := range errorChan {
			fmt.Println("error in leader-election process:", err)
		}
	}()

	// Block forever.
	select {}
}
