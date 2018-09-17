package main

import (
	"github.com/fighterlyt/test/grpc/add"
	"context"
	"net"
	"log"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	add.RegisterAddServer(grpcServer, addServer{})
	grpcServer.Serve(lis)
}

type addServer struct {
}

func (addServer) Add(ctx context.Context, data *add.Data) (*add.Data, error) {
	data.Value++
	return data, nil
}
