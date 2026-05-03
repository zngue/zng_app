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
)

// HttpServer 封装 http.Server，提供优雅启动和关闭功能
type HttpServer struct {
	*http.Server
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

// Stop 优雅关闭 HTTP 服务器
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

// NewHttpServer 创建新的 HttpServer 实例
func NewHttpServer(addr string, handler http.Handler) *HttpServer {
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return &HttpServer{Server: server}
}
