package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zngue/zng_app/app/server/middleware"
)

type HttpServer struct {
	*http.Server
}

type HttpServerOption func(*HttpServer)

func WithMiddlewareRegistry(r *middleware.Registry) HttpServerOption {
	return func(s *HttpServer) {
		middleware.SetRegistry(r)
	}
}

func (s *HttpServer) Start() error {
	go func() {
		log.Printf("HTTP server starting on %s", s.Server.Addr)
		if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
			panic(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down HTTP server...")
	if err := s.Stop(); err != nil {
		return fmt.Errorf("failed to gracefully shutdown HTTP server: %w", err)
	}

	log.Println("HTTP server stopped")
	return nil
}
func (s *HttpServer) Stop() error {
	if s.Server == nil {
		return fmt.Errorf("server is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
		return err
	}

	return nil
}

func NewHttpServer(addr string, handler http.Handler, opts ...HttpServerOption) *HttpServer {
	s := &HttpServer{}
	for _, opt := range opts {
		opt(s)
	}
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	s.Server = server
	return s
}
