package main

import (
	"context"
	"github.com/fighterlyt/test/grpc/add"
	"google.golang.org/grpc"
	"log"
	"net"
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
	log.Println("add", data)
	data.Value++
	return data, nil
}
