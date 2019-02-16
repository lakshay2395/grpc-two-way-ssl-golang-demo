package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/net/context"

	"../api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// //Authentication - structure with username and password
// type Authentication struct {
// 	Login    string
// 	Password string
// }

// //BELOW TWO METHODS SHOULD BE IMPLEMENTED AS PER REQUIREMENT FOR PER RPC METADATA

// // GetRequestMetadata - gets the current request metadata
// func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
// 	return map[string]string{
// 		"login":    a.Login,
// 		"password": a.Password,
// 	}, nil
// }

// // RequireTransportSecurity - indicates whether the credentials requires transport security
// func (a *Authentication) RequireTransportSecurity() bool {
// 	return true
// }

func main() {

	var conn *grpc.ClientConn
	trustedCerts := []string{"cert/server.crt"}
	clientIdentity, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if err != nil {
		fmt.Println(fmt.Errorf("failed to create x509 keypair: %v", err))
	}
	certPool := x509.NewCertPool()
	for _, cert := range trustedCerts {
		certFile, err := ioutil.ReadFile(cert)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to read cert file : %v", err))
		}
		certPool.AppendCertsFromPEM(certFile)
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clientIdentity},
		RootCAs:      certPool,
		ServerName:   "Sample",
	})
	conn, err = grpc.Dial(":7777", grpc.WithTransportCredentials(creds))
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
