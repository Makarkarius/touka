package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type server struct {
	cfg    Config
	server *http.Server

	requestQueue chan RequestBatch

	requester Requester
	storage   Storager

	logger *zap.Logger

	rabbitConn    *amqp.Connection
	rabbitChannel *amqp.Channel
}

func NewServer(cfg Config) (*server, error) {
	logger, err := cfg.LoggerCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	requester, err := cfg.RequesterCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build requester: %w", err)
	}
	requester.UseLogger(logger)

	storage, err := cfg.StorageCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build storage: %w", err)
	}
	storage.UseLogger(logger)

	rabbitConn, err := amqp.Dial(cfg.RabbitURI)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rabbit: %w", err)
	}
	channel, err := rabbitConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open rabbit channel: %w", err)
	}

	s := &server{
		cfg: cfg,
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:      nil,
			ReadTimeout:  time.Duration(cfg.ReadTimeoutSec) * time.Second,
			WriteTimeout: time.Duration(cfg.WriteTimeoutSec) * time.Second,
		},
		requestQueue:  make(chan RequestBatch, cfg.RequestQueueSize),
		requester:     requester,
		storage:       storage,
		logger:        logger,
		rabbitConn:    rabbitConn,
		rabbitChannel: channel,
	}

	router := mux.NewRouter()
	router.HandleFunc("/get", s.handle(handleGet))
	s.server.Handler = router

	return s, nil
}

// handling only one request batch at one time
func (s *server) serveHandlers(ctx context.Context) {
	s.logger.Info("start serving handlers")
	defer s.logger.Info("stop serving handlers")

	wg := sync.WaitGroup{}
	for {
		select {
		case batch := <-s.requestQueue:
			responseQueue := make(chan *Response, s.cfg.ResponseQueueSize)
			successful := 0

			wg.Add(2)
			go func() {
				defer wg.Done()
				successful = s.requester.GetBatch(ctx, responseQueue, batch)
			}()
			go func() {
				defer wg.Done()
				s.storage.InsertChan(ctx, responseQueue)
			}()
			wg.Wait()

			msg, err := formRabbitMsg(batch.CmdCode, successful)
			if err != nil {
				s.logger.Error("rabbit form message error", zap.Error(err))
				continue
			}
			if err := s.rabbitPublish(ctx, msg); err != nil {
				s.logger.Error("rabbit publish message error", zap.Error(err))
				continue
			}
			s.logger.Info("rabbit successful message publish", zap.Int("successfulCount", successful))
		case <-ctx.Done():
			return
		}
	}
}

func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.serveHandlers(ctx)
	}()

	defer func() {
		_ = s.logger.Sync()
		_ = s.rabbitConn.Close()
		cancel()
		wg.Wait()
		s.logger.Info("server stopped")
	}()

	s.logger.Info("starting server",
		zap.String("host", s.cfg.Host),
		zap.Int("port", s.cfg.Port))
	return s.server.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *server) handle(f func(s *server, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(s, w, r)
	}
}
