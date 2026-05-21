package p2p

// HandshakerFunc is ... ?
type HandshakerFunc func(Peer) error

func NOPHandshakeFunc(Peer) error { return nil }
