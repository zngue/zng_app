package rpc

import (
	"google.golang.org/grpc"
	"log"
)

type Rpc struct {
}

func (r *Rpc) Start() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
}
func (r *Rpc) Stop() {

}
