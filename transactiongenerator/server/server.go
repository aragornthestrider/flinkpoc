package server

import (
	"context"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	ctx             context.Context
	apiReadTimeout  int
	apiWriteTimeout int
	IsHealthy       *atomic.Value
	IsReady         *atomic.Value
	srv             *http.Server
	Logger          *zap.Logger
}

func NewServer(ctx context.Context, port, apiReadTimeout, apiWriteTimeout int, logger *zap.Logger) *Server {
	server := new(Server)
	server.IsHealthy = &atomic.Value{}
	server.IsReady = &atomic.Value{}
	server.srv = &http.Server{
		Handler:      server.NewRouter(),
		Addr:         ":" + strconv.Itoa(port),
		ReadTimeout:  time.Duration(apiReadTimeout) * time.Second,
		WriteTimeout: time.Duration(apiWriteTimeout) * time.Second,
	}
	server.Logger = logger
	server.ctx = ctx
	return server
}

func (s *Server) Context() context.Context {
	return s.ctx
}

func (s *Server) SetHealthy(status bool) {
	s.IsHealthy.Store(status)
}

func (s *Server) SetReady(status bool) {
	s.IsReady.Store(status)
}

func (s *Server) Run() error {
	errCh := make(chan error)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-s.Context().Done():
		s.Logger.Debug("Main context finished")
		return s.Stop()
	case err := <-errCh:
		s.Logger.Error("Error received from server", zap.Error(err))
		return s.Stop()
	}
}

func (s *Server) Stop() error {
	if err := s.srv.Shutdown(s.Context()); err != nil {
		return err
	}
	s.Logger.Info("Server shut down successfully")
	return nil
}
