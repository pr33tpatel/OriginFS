package p2p

import "net"

// RPC represents any arbitrary data being sent over
// each transport between two nodes
type RPC struct {
	From    net.Addr
	Payload []byte
}
