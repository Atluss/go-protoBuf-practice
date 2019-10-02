package main

import (
	"context"
	"fmt"
	"github.com/Atluss/protoBufPractice/pkg/cctx"
	"github.com/Atluss/protoBufPractice/pkg/v1/proto/streaming"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	addr := fmt.Sprintf(":%s", "14534")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(cctx.SignalContext(context.Background(), "gRPC streaming example"))
	defer cancel()

	grpcServer := grpc.NewServer()

	stream := streaming.New()
	streaming.RegisterStreamingServer(grpcServer, stream)

	go func() {
		log.Printf("Start streaming server")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to start server: %s", err)
		}
		cancel()
	}()

	<-ctx.Done()
	log.Printf("Stop streaming server")
}
