package p2p

import "errors"

// ErrInvalidHandshake is returned if the handshake between the local and remote node cannot be established
var ErrInvalidHandshake = errors.New("invalid handshake")

// HandshakerFunc is ... ?
type HandshakerFunc func(Peer) error

func NOPHandshakeFunc(Peer) error { return nil }
