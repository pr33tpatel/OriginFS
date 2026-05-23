package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pr33tpatel/OriginFS/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	// return fmt.Errorf("failed the onpeer func")
	fmt.Println("doing some logic with peer outside of TCPTransport")
	return nil
}

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO: OnPeer function
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageOrigin:     ":3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}
	s := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
