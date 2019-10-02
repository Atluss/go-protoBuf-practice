package streaming

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"sync"
)

const TokenHeader = "access-token"

type ServerPingPong struct {
	Clients    map[string]chan ServerMsg
	streamsMtx sync.RWMutex
}

func New() *ServerPingPong {
	return &ServerPingPong{
		Clients:    map[string]chan ServerMsg{},
		streamsMtx: sync.RWMutex{},
	}
}

func (s *ServerPingPong) Stream(srv Streaming_StreamServer) (err error) {
	tkn, ok := s.extractToken(srv.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "missing token header")
	}

	go s.send(srv, tkn)
	go s.receive(srv, tkn)

	<-srv.Context().Done()
	return srv.Context().Err()
}

func (s *ServerPingPong) receive(srv Streaming_StreamServer, tkn string) error {
	for {

		res, err := srv.Recv()

		if s, ok := status.FromError(err); ok && s.Code() == codes.Canceled {
			log.Printf("client disconnected by himself: %s", s.Message())
			return nil
		} else if err == io.EOF {
			log.Printf("client disconnected EOF")
			return nil
		} else if err != nil {
			log.Printf("failed to receive a note : %v", err)
			return err
		}

		switch evt := res.Event.(type) {
		case *ClientMsg_ClnSendPingPong:

			msg := &SendPingPong{Msg: "Server say!"}

			log.Printf("msg: %+v", res)

			if res.GetClnSendPingPong().Type == PONG_NUM_PING {
				log.Printf("%s send PING", tkn)
				msg.Type = PONG_NUM_PONG
			} else {
				log.Printf("%s send PONG", tkn)
				msg.Type = PONG_NUM_PING
			}

			s.Clients[tkn] <- ServerMsg{
				Timestamp: ptypes.TimestampNow(),
				Event: &ServerMsg_SrvSendPingPong{
					SrvSendPingPong: msg,
				},
			}
		case *ClientMsg_ClnShutdown_:
			log.Printf("Client close connection")
		default:
			log.Printf("(bot connection) unexpected event from the server: %T", evt)
			return nil
		}
	}
}

func (s *ServerPingPong) send(srv Streaming_StreamServer, tkn string) {
	stream := s.openStream(tkn)
	defer s.closeStream(tkn)

	for {
		select {
		case <-srv.Context().Done():
			return
		case res := <-stream:
			if s, ok := status.FromError(srv.Send(&res)); ok {
				switch s.Code() {
				case codes.OK:
				case codes.Unavailable, codes.Canceled, codes.DeadlineExceeded:
					log.Printf("client (%s) terminated connection", tkn)
					return
				default:
					log.Printf("error send message client: %s", tkn)
					return
				}
			}
		}
	}
}

func (s *ServerPingPong) openStream(tkn string) (stream chan ServerMsg) {
	stream = make(chan ServerMsg, 100)
	s.streamsMtx.Lock()
	s.Clients[tkn] = stream
	s.streamsMtx.Unlock()

	log.Printf("client connected: %s", tkn)
	return stream
}

func (s *ServerPingPong) closeStream(tkn string) {
	s.streamsMtx.Lock()
	if stream, ok := s.Clients[tkn]; ok {
		close(stream)
		delete(s.Clients, tkn)
	}
	s.streamsMtx.Unlock()
	log.Printf("client disconnected: %s", tkn)
}

func (s *ServerPingPong) extractToken(ctx context.Context) (tkn string, ok bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md[TokenHeader]) == 0 {
		return "", false
	}

	return md[TokenHeader][0], true
}
