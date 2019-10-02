package client_streaming

import (
	"context"
	"fmt"
	client "github.com/Atluss/protoBufPractice/pkg/v1/proto/streaming"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

type ClientStream struct {
	client.StreamingClient
	Host, Token string
	Message     chan client.ClientMsg
}

func Client(host, token string) *ClientStream {
	return &ClientStream{
		Host:    host,
		Token:   token,
		Message: make(chan client.ClientMsg, 10),
	}
}

func (c *ClientStream) Run(ctx context.Context) error {
	connCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(connCtx, c.Host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Connection close error: %s", err)
		}
	}()

	c.StreamingClient = client.NewStreamingClient(conn)

	err = c.stream(ctx)

	return fmt.Errorf("error stream")
}

func (c *ClientStream) stream(ctx context.Context) error {
	md := metadata.New(map[string]string{client.TokenHeader: c.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := c.StreamingClient.Stream(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := stream.CloseSend(); err != nil {
			log.Printf("Connection close error: %s", err)
		}

	}()

	go c.send(stream)
	return c.receive(stream)
}

func (c *ClientStream) send(client client.Streaming_StreamClient) {
	for {
		select {
		case <-client.Context().Done():
			log.Printf("client send loop disconnected")
			return
		case data := <-c.Message:
			//log.Printf("%+v", data)

			if err := client.Send(&data); err != nil {
				log.Fatalf("Failed to send a note: %v", err)
			}
		}
		log.Printf("message printed")
	}
}

func (c *ClientStream) receive(cnt client.Streaming_StreamClient) error {
	for {

		log.Printf("Wait message from server")

		res, err := cnt.Recv()

		if s, ok := status.FromError(err); ok && s.Code() == codes.Canceled {
			log.Println("stream canceled (usually indicates shutdown)")
			return nil
		} else if err == io.EOF {
			log.Println("stream closed by server")
			return nil
		} else if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
			return err
		}

		switch evt := res.Event.(type) {
		case *client.ServerMsg_SrvSendPingPong:
			pingPong := res.GetSrvSendPingPong()

			msg := &client.SendPingPong{Msg: "send srv msg"}

			if pingPong.Type == client.PONG_NUM_PING {
				msg.Type = client.PONG_NUM_PONG
				log.Printf("Server semd to me %s PING", c.Token)
			} else {
				msg.Type = client.PONG_NUM_PING
				log.Printf("Server semd to me %s PONG", c.Token)
			}

			time.Sleep(5 * time.Second)

			c.Message <- client.ClientMsg{
				Timestamp: ptypes.TimestampNow(),
				Event: &client.ClientMsg_ClnSendPingPong{
					ClnSendPingPong: msg},
			}

		default:
			log.Printf("unexpected event from the server: %T", evt)
			return nil
		}
	}
}
