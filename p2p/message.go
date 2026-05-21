package p2p

// Message represents any arbitrary data being sent over
// each transport between two nodes
type Message struct {
	Payload []byte
}
