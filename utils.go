package main

import (
	"fmt"
	"strconv"
	"strings"
)

// parseSequence returns the sequence number of the given Zookeeper sequential node,
// assuming the number starts immediately after the given prefix.
func parseSequence(sequentialNodePath, prefix string) (int, error) {
	sequenceString := strings.TrimPrefix(sequentialNodePath, prefix)
	sequence, err := strconv.ParseInt(sequenceString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error in strconv.ParseInt call: %w", err)
	}

	return int(sequence), nil
}
