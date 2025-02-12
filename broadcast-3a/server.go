package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"sync"
)

type Server struct {
	node *maelstrom.Node

	idsMu sync.RWMutex
	ids   []int

	topologyMu sync.RWMutex
	Topology   map[string][]string
}

type BroadcastMessage struct {
	Type    string `json:"type"`
	Message int    `json:"message"`
}

type BroadcastMessageResponse struct {
	Type string `json:"type"`
}

type ReadMessage struct {
	Type string `json:"type"`
}

type ReadMessageResponse struct {
	Type     string `json:"type"`
	Messages []int  `json:"messages"`
}

type TopologyMessage struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

type TopologyMessageResponse struct {
	Type string `json:"type"`
}

func NewServer(n *maelstrom.Node) *Server {
	return &Server{node: n}
}

func (s *Server) ReadTopologyHandler(msg maelstrom.Message) error {
	defer s.topologyMu.Unlock()

	var body TopologyMessage
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.topologyMu.Lock()
	s.Topology = body.Topology

	topologyMessageResponse := TopologyMessageResponse{
		Type: "topology_ok",
	}

	return s.node.Reply(msg, topologyMessageResponse)
}

func (s *Server) BroadcastHandler(msg maelstrom.Message) error {
	defer s.idsMu.Unlock()

	var body BroadcastMessage
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.idsMu.Lock()
	s.ids = append(s.ids, body.Message)

	broadcastMessageResponse := BroadcastMessageResponse{Type: "broadcast_ok"}

	return s.node.Reply(msg, broadcastMessageResponse)
}

func (s *Server) ReadHandler(msg maelstrom.Message) error {
	defer s.idsMu.Unlock()

	var body ReadMessage
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.idsMu.Lock()
	readMessageResponse := ReadMessageResponse{Type: "read_ok", Messages: s.ids}

	return s.node.Reply(msg, readMessageResponse)
}
