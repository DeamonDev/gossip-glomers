package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"go.uber.org/zap"
)

type Server struct {
	node *maelstrom.Node
	kv   *maelstrom.KV

	nodeID string
	peers  []string

	logger *zap.Logger
}

type AddMessage struct {
	Type  string `json:"type"`
	Delta int    `json:"delta"`
}

type AddMessageResponse struct {
	Type string `json:"type"`
}

type ReadMessage struct {
	Type string `json:"type"`
}

type ReadMessageResponse struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type InitMessageResponse struct {
	Type string `json:"type"`
}

func NewServer(n *maelstrom.Node, l *zap.Logger) *Server {
	kv := maelstrom.NewKV("counter", n)
	s := &Server{node: n, kv: kv, logger: l}

	n.Handle("init", s.initHandler)

	return s
}

func (s *Server) initHandler(msg maelstrom.Message) error {
	var body maelstrom.InitMessageBody
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.logger.Info("init message received", zap.String("node_id", body.NodeID))

	s.nodeID = body.NodeID

	for _, id := range body.NodeIDs {
		if id != s.nodeID {
			s.peers = append(s.peers, id)
		}
	}

	s.logger.Info("peers discovered due to initial message: ", zap.Strings("peers", s.peers))

	initMessageResponse := InitMessageResponse{Type: "init_ok"}

	return s.node.Reply(msg, initMessageResponse)
}

func (s *Server) Run() error {
	return s.node.Run()
}
