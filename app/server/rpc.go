package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zngue/zng_app/app/server/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	*grpc.Server
	listener net.Listener
}

type GRPCServerOption func(*grpcServerConfig)

type grpcServerConfig struct {
	grpcOpts   []grpc.ServerOption
	registry   *middleware.Registry
	reflection bool
}

func WithGRPCOption(opts ...grpc.ServerOption) GRPCServerOption {
	return func(c *grpcServerConfig) {
		c.grpcOpts = append(c.grpcOpts, opts...)
	}
}

func WithGRPCMiddlewareRegistry(r *middleware.Registry) GRPCServerOption {
	return func(c *grpcServerConfig) {
		c.registry = r
	}
}

func WithGRPCReflection() GRPCServerOption {
	return func(c *grpcServerConfig) {
		c.reflection = true
	}
}

func unaryServerInterceptor(r *middleware.Registry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		operation := info.FullMethod
		h := r.Build(operation, func(ctx context.Context) (any, error) {
			return handler(ctx, req)
		})
		out, err := h(ctx)
		afterH := r.BuildAfter(operation, func(ctx context.Context, err error, in any, rs any) {})
		afterH(ctx, err, req, out)
		return out, err
	}
}

func (s *GRPCServer) Start() error {
	go func() {
		log.Printf("gRPC server starting on %s", s.listener.Addr().String())
		if err := s.Serve(s.listener); err != nil && err != grpc.ErrServerStopped {
			log.Printf("gRPC server error: %v", err)
			panic(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	if err := s.Stop(); err != nil {
		return fmt.Errorf("failed to gracefully shutdown gRPC server: %w", err)
	}
	log.Println("gRPC server stopped")
	return nil
}

func (s *GRPCServer) Stop() error {
	if s.Server == nil {
		return fmt.Errorf("server is nil")
	}
	done := make(chan struct{})
	go func() {
		s.Server.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-time.After(5 * time.Second):
		s.Server.Stop()
		return fmt.Errorf("gRPC server graceful shutdown timed out, forced stop")
	}
}

func NewGRPCServer(addr string, opts ...GRPCServerOption) (*GRPCServer, error) {
	cfg := &grpcServerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.registry != nil {
		cfg.grpcOpts = append(cfg.grpcOpts, grpc.UnaryInterceptor(unaryServerInterceptor(cfg.registry)))
		middleware.SetRegistry(cfg.registry)
	}
	grpcServer := grpc.NewServer(cfg.grpcOpts...)
	if cfg.reflection {
		reflection.Register(grpcServer)
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", addr, err)
	}
	return &GRPCServer{
		Server:   grpcServer,
		listener: listener,
	}, nil
}

func (s *GRPCServer) GetAddr() string {
	if s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}
