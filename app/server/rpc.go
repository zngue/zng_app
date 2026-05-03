package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServer 封装 gRPC 服务器，提供优雅启动和关闭功能
type GRPCServer struct {
	*grpc.Server
	listener net.Listener
}

// Start 启动 gRPC 服务器并处理优雅关闭
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
	s.Server.Stop()
	log.Println("gRPC server stopped")
	return nil
}
func (s *GRPCServer) Stop() error {
	if s.Server == nil {
		log.Println("gRPC server is nil, nothing to stop")
		return nil
	}
	s.Server.GracefulStop()
	return nil
}
func NewGRPCServer(addr string, opts ...grpc.ServerOption) (server *GRPCServer, err error) {
	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer(opts...)
	// 创建监听器
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		err = fmt.Errorf("failed to listen on address %s: %w", addr, err)
		return
	}
	server = &GRPCServer{
		Server:   grpcServer,
		listener: listener,
	}
	return server, nil
}

// NewGRPCServerWithReflection 创建带反射功能的 gRPC 服务器实例
func NewGRPCServerWithReflection(addr string, opts ...grpc.ServerOption) (*GRPCServer, error) {
	server, err := NewGRPCServer(addr, opts...)
	if err != nil {
		return nil, err
	}
	// 注册反射服务，便于调试和测试
	reflection.Register(server.Server)
	return server, nil
}

// GetAddr 返回服务器监听地址
func (s *GRPCServer) GetAddr() string {
	if s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}
