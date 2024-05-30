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

	// Start participating in the election.
	electedChan, errorChan := zook.Participate()

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
