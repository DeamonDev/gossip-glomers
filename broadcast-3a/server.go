package main

import (
	"encoding/json"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
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
	s := &Server{node: n}

	n.Handle("topology", s.readTopologyHandler)
	n.Handle("broadcast", s.broadcastHandler)
	n.Handle("read", s.readHandler)

	return s
}

func (s *Server) Run() error {
	return s.node.Run()
}

func (s *Server) readTopologyHandler(msg maelstrom.Message) error {
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

func (s *Server) broadcastHandler(msg maelstrom.Message) error {
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

func (s *Server) readHandler(msg maelstrom.Message) error {
	defer s.idsMu.Unlock()

	var body ReadMessage
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.idsMu.Lock()
	readMessageResponse := ReadMessageResponse{Type: "read_ok", Messages: s.ids}

	return s.node.Reply(msg, readMessageResponse)
}
