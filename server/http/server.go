package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var ErrServiceNotReady = errors.New("http service: not started yet")

type Config struct {
	Port    int `mapstructure:"PORT"  default:"8080"`
	Timeout int `mapstructure:"TIMEOUT"  default:"60"`
}

type Server struct {
	http     *http.Server
	listener net.Listener
	ready    bool
	runErr   error
}

func New(cfg Config) *Server {
	httpSrv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		ReadHeaderTimeout: time.Second * time.Duration(cfg.Timeout),
	}

	server := &Server{
		http:   httpSrv,
		runErr: nil,
	}

	return server
}

func (s *Server) Status() error {
	if !s.ready {
		return ErrServiceNotReady
	}
	if s.runErr != nil {
		return s.runErr
	}

	return nil
}

func (s *Server) Serve() error {
	netListen, err := net.Listen("tcp", s.http.Addr)
	if err != nil {
		return fmt.Errorf("faile to create NewWssServer network con: %w", err)
	}
	log.Info(log.WithField("NewHttpServer server at", s.http.Addr))
	s.listener = netListen
	s.ready = true

	err = s.http.Serve(s.listener)
	s.runErr = err

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.http.Shutdown(ctx)
	s.runErr = err

	return err
}
