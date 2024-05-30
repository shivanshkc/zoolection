package election

import (
	"fmt"
	"strconv"
	"strings"
)

// parseESNodeSequence parses the sequence number of an ephemeral-sequential zNode.
func parseESNodeSequence(nodePath string) (int64, error) {
	// Get the name of the ES node. Example: /election/candidate -> candidate
	esNodeName := ephemeralNodePath[strings.LastIndex(ephemeralNodePath, "/")+1:]
	// Get the index at which the sequence number begins.
	seqStartIndex := strings.LastIndex(nodePath, esNodeName) + len(esNodeName)

	// Parse the sequence.
	value, err := strconv.ParseInt(nodePath[seqStartIndex:], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error in strconv.ParseInt call: %w", err)
	}

	return value, nil
}

// mapESNodeSequence parses all given node paths to get the sequence number and returns a map of path to sequence.
func mapESNodeSequence(nodePaths []string) (map[string]int64, error) {
	// Create the map to be returned.
	pathSeqMap := make(map[string]int64, len(nodePaths))

	// Parse all paths.
	for _, path := range nodePaths {
		seq, err := parseESNodeSequence(path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sequence for node %s: %w", path, err)
		}
		// Add to map.
		pathSeqMap[path] = seq
	}

	return pathSeqMap, nil
}
