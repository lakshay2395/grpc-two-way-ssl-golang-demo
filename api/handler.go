package api

import (
	"log"

	"golang.org/x/net/context"
)

//Server - server represents grpc server
type Server struct {
}

//SayHello - sayhello handler function for grpc
func (s *Server) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Message Received : %s", in.Greeting)
	return &PingMessage{Greeting: "bar"}, nil
}
