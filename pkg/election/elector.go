package election

// Elector represents the minimal functionality required to implement leader election.
type Elector interface {
	// Init does the initial setup required to participate in the election.
	//
	// Leave empty if the implementation doesn't require any setup.
	Init() error

	// Participate in the election.
	//
	// This method returns two read-only channels.
	// If the node is elected as the leader, the first channel receives an event and the participation ends.
	// If there's any error in the election process, the second channel reports it.
	Participate() (<-chan struct{}, <-chan error)
}
