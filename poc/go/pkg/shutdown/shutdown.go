package shutdown

// Please add the dependencies if you add your own priority here.
// Otherwise investigating deadlocks at shutdown is much more complicated.

const (
	PriorityP2PManager = iota // no dependencies
	PriorityKademliaDHT
	PriorityPeerDiscovery
)
