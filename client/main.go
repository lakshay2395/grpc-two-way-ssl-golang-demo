package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	"../api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//Authentication - structure with username and password
type Authentication struct {
	Login    string
	Password string
}

//BELOW TWO METHODS SHOULD BE IMPLEMENTED AS PER REQUIREMENT FOR PER RPC METADATA

// GetRequestMetadata - gets the current request metadata
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil
}

// RequireTransportSecurity - indicates whether the credentials requires transport security
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

func main() {

	var conn *grpc.ClientConn

	creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "Sample")
	if err != nil {
		log.Fatal(err)
	}

	auth := Authentication{
		Login:    "lakshay",
		Password: "sample",
	}

	conn, err = grpc.Dial(":7777", grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := api.NewPingClient(conn)
	response, err := client.SayHello(context.Background(), &api.PingMessage{Greeting: "foo"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Greeting)
}
