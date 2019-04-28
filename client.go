package main

import (
	"context"
	greeting "github.com/Atluss/protoBufPractice/proto/proto"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := greeting.NewGreetingClient(conn)

	// Contact the server and print out its response.
	name := "client"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &greeting.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Msg: %s", r.Message)
}
