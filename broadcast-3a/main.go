package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
)

func main() {
	n := maelstrom.NewNode()
	s := NewServer(n)

	n.Handle("topology", s.ReadTopologyHandler)
	n.Handle("broadcast", s.BroadcastHandler)
	n.Handle("read", s.ReadHandler)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
