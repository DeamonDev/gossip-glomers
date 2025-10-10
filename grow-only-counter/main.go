package main

import (
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"go.uber.org/zap"
)

func main() {
	n := maelstrom.NewNode()

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"/tmp/node.log"}
	cfg.ErrorOutputPaths = []string{"/tmp/node_error.log"}

	logger, _ := cfg.Build()
	defer logger.Sync()

	s := NewServer(n, logger)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

}
