package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

type LocalReadMessageResponse struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

func NewServer(n *maelstrom.Node) *Server {
	kv := maelstrom.NewSeqKV(n)
	s := &Server{node: n, kv: kv}

	n.Handle("init", s.initHandler)
	n.Handle("local_read", s.localReadHandler)

	return s
}

func (s *Server) initHandler(msg maelstrom.Message) error {
	var body maelstrom.InitMessageBody
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.nodeID = body.NodeID

	err := initLogger(s)
	if err != nil {
		return err
	}

	s.logger.Info("init message received", zap.String("node_id", body.NodeID))

	for _, id := range body.NodeIDs {
		if id != s.nodeID {
			s.peers = append(s.peers, id)
		}
	}

	s.logger.Info("peers discovered due to initial message:", zap.Strings("peers", s.peers))

	initMessageResponse := InitMessageResponse{Type: "init_ok"}

	return s.node.Reply(msg, initMessageResponse)
}

func (s *Server) localReadHandler(msg maelstrom.Message) error {
	v, err := s.kv.ReadInt(context.Background(), "counter")
	if err != nil {
		return err
	}

	localReadMessageResponse := LocalReadMessageResponse{Type: "local_read_ok", Value: v}

	return s.node.Reply(msg, localReadMessageResponse)
}

func initLogger(s *Server) error {
	err := os.MkdirAll("/tmp/logs", 0755)
	if err != nil {
		return err
	}

	logFile := fmt.Sprintf("/tmp/logs/node_%s.log", s.nodeID)
	errLogFile := fmt.Sprintf("/tmp/logs/node_%s.err.log", s.nodeID)

	f, _ := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	ef, _ := os.OpenFile(errLogFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)

	cfg := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(f),
		zap.DebugLevel,
	)

	errCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(ef),
		zap.ErrorLevel,
	)

	go initCache(s)

	s.logger = zap.New(zapcore.NewTee(core, errCore))
	return nil
}

func initCache(s *Server) {
	err := s.kv.Write(context.Background(), "counter", 0)
	if err != nil {
		s.logger.Fatal("failed to initialize counter", zap.Error(err))
	} else {
		s.logger.Info("initialized counter for myself")
	}

	for _, peerID := range s.peers {
		err = s.kv.Write(context.Background(), peerID, 0)
		if err != nil {
			s.logger.Fatal("failed to initialize counter for", zap.String("peerID", peerID), zap.Error(err))
		} else {
			s.logger.Info("initialized counter for", zap.String("peerID", peerID))
		}
	}
}

func (s *Server) Run() error {
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(s.logger)

	return s.node.Run()
}
