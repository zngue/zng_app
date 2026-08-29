package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	*http.Server
	shutdownTimeout time.Duration
	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
}

type HTTPOption func(*HTTPServer)

func NewHTTPServer(addr string, handler http.Handler, opts ...HTTPOption) *HTTPServer {
	s := &HTTPServer{
		shutdownTimeout: 10 * time.Second,
		readTimeout:     10 * time.Second,
		writeTimeout:    10 * time.Second,
		idleTimeout:     30 * time.Second,
	}
	for _, opt := range opts {
		opt(s)
	}
	s.Server = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
	}
	return s
}

func WithTimeouts(read, write, idle time.Duration) HTTPOption {
	return func(s *HTTPServer) {
		s.readTimeout = read
		s.writeTimeout = write
		s.idleTimeout = idle
	}
}

func WithShutdownTimeout(d time.Duration) HTTPOption {
	return func(s *HTTPServer) {
		s.shutdownTimeout = d
	}
}

func (s *HTTPServer) Start() error {
	log.Printf("http server starting on %s", s.Addr)
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *HTTPServer) Stop() error {
	if s.Server == nil {
		return fmt.Errorf("server is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.Shutdown(ctx)
}
