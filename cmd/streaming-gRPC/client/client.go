package main

import (
	"context"
	"github.com/Atluss/protoBufPractice/cmd/streaming-gRPC/client/client_streaming"
	"github.com/Atluss/protoBufPractice/pkg/cctx"
	"github.com/Atluss/protoBufPractice/pkg/v1/proto/streaming"
	"github.com/golang/protobuf/ptypes"
	"log"
)

var orchestraAddress = "localhost:14534"

var Clients = []string{
	"CLIENT-1",
	//"CLIENT-2",
	//"CLIENT-3",
	//"CLIENT-4",
}

func main() {
	systemName := "BotsConn"

	ctx, cancel := context.WithCancel(cctx.SignalContext(context.Background(), systemName))
	defer cancel()

	RunClients()

	<-ctx.Done()
	log.Printf("%s stop", systemName)
}

func RunClients() {
	for i := 0; i < len(Clients); i++ {
		// test registration request and registration status
		log.Printf("Connection bot: %s", Clients[i])
		go clientsConnection(Clients[i])
	}
}

func clientsConnection(tkn string) {

	ctx, cancel := context.WithCancel(cctx.SignalContext(context.Background(), "Streaming client"))
	defer cancel()

	client := client_streaming.Client(orchestraAddress, tkn)

	client.Message <- streaming.ClientMsg{
		Timestamp: ptypes.TimestampNow(),
		Event: &streaming.ClientMsg_ClnSendPingPong{
			ClnSendPingPong: &streaming.SendPingPong{
				Type: streaming.PONG_NUM_PONG,
				Msg:  "send msg to server",
			},
		},
	}

	if err := client.Run(ctx); err != nil {
		log.Printf("error stream client: %s", err)
		return
	}

}
