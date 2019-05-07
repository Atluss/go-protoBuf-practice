package main

import (
	"context"
	"github.com/Atluss/protoBufPractice/pkg/v1/proto/greeting"
	"google.golang.org/grpc"
	"log"
	"net"
)

// generate proto file, with grpc plugin to request and answer
//go:generate protoc -I ../../pkg/v1/proto/greeting --go_out=plugins=grpc:../../pkg/v1/proto/greeting ../../pkg/v1/proto/greeting/greeting.proto

// server is used to implement greeting.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *greeting.HelloRequest) (*greeting.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &greeting.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greeting.RegisterGreetingServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
