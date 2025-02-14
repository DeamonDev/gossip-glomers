package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
)

func main() {
	n := maelstrom.NewNode()
	s := NewServer(n)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
