package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over an established TCP connection
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn

	// if we dial and retrieve a connection => outbound == true
	// if we accept and retrieve a connection => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	shakeHands    HandshakerFunc
	decoder       Decoder

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

// NewTCPTransport initializes a TCP connection with a listen address and returns the TCP connection
func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		shakeHands:    NOPHandshakeFunc,
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		fmt.Printf("new incoming connection: %+v\n", conn)
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true) // NOTE: should be "false" i think

	_ = peer

	if err := t.shakeHands(conn); err != nil {
		// if there is an error in the handshake, then we need to drop the connection
		conn.Close()
		return
	}

	// read loop
	msg := &Temp{}
	for {
		// n, _ := conn.Read(bud)
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
		}
		// msg := buf[:n]
	}
}
