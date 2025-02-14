package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"slices"
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
	s := &Server{node: n}

	n.Handle("topology", s.readTopologyHandler)
	n.Handle("broadcast", s.broadcastHandler)
	n.Handle("broadcast_ok", s.broadcastOkHandler)
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

func (s *Server) broadcastOkHandler(msg maelstrom.Message) error {
	return nil
}

func (s *Server) broadcastHandler(msg maelstrom.Message) error {
	var body BroadcastMessage
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.idsMu.Lock()

	//Here I've assumed broadcast messages are unique!
	if slices.Contains(s.ids, body.Message) {
		s.idsMu.Unlock()
		return nil
	}

	s.ids = append(s.ids, body.Message)
	s.idsMu.Unlock()

	wg := &sync.WaitGroup{}

	peers := s.Topology[s.node.ID()]

	wg.Add(len(peers))
	for _, peerNodeId := range peers {

		if s.node.ID() == msg.Src || peerNodeId == msg.Src {
			wg.Done()
			continue
		}

		go func(id string) {
			defer wg.Done()
			if err := s.node.Send(peerNodeId, body); err != nil {
				panic(err)
			}
		}(peerNodeId)
	}

	wg.Wait()

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
