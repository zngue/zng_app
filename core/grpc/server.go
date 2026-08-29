package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	*grpc.Server
	addr string
}

type Option func(*GRPCServer)

func NewGRPCServer(addr string, opts ...grpc.ServerOption) *GRPCServer {
	srv := grpc.NewServer(opts...)
	return &GRPCServer{
		Server: srv,
		addr:   addr,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Printf("grpc server starting on %s", s.addr)
	return s.Serve(lis)
}

func (s *GRPCServer) Stop() error {
	s.GracefulStop()
	return nil
}

func (s *GRPCServer) Addr() string {
	return s.addr
}
