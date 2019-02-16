package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"../api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// private type for Context keys
type contextKey int

const (
	clientIDKey contextKey = iota
)

// func credMatcher(headerName string) (mdName string, ok bool) {
// 	if headerName == "Login" || headerName == "Password" {
// 		return headerName, true
// 	}
// 	return "", false
// }

// func authenticateClient(ctx context.Context, s *api.Server) (string, error) {
// 	if md, ok := metadata.FromIncomingContext(ctx); ok {
// 		login := strings.Join(md["login"], "")
// 		password := strings.Join(md["password"], "")
// 		if login != "lakshay" {
// 			return "", fmt.Errorf("unknown user %s", login)
// 		}
// 		if password != "sample" {
// 			return "", fmt.Errorf("bad password %s", password)
// 		}
// 		log.Printf("authenticated client: %s", login)
// 		return "42", nil
// 	}
// 	return "", fmt.Errorf("missing credentials")
// }

func startGRPCServer(address, certFile string, keyFile string, trustedCerts []string) error {
	serverIdentity, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to create x509 keypair: %v", err)
	}
	certPool := x509.NewCertPool()
	for _, cert := range trustedCerts {
		certFile, err := ioutil.ReadFile(cert)
		if err != nil {
			return fmt.Errorf("failed to read cert file : %v", err)
		}
		certPool.AppendCertsFromPEM(certFile)
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := api.Server{}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverIdentity},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})
	opts := []grpc.ServerOption{grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)
	api.RegisterPingServer(grpcServer, &s)
	log.Printf("started HTTP/2 2way ssl gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}
	return nil
}

// func startRESTServer(address, grpcAddress, certFile string) error {
// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()
// 	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(credMatcher))
// 	creds, err := credentials.NewClientTLSFromFile(certFile, "sample")
// 	if err != nil {
// 		return fmt.Errorf("could not load TLS certificate: %s", err)
// 	}
// 	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
// 	err = api.RegisterPingHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
// 	if err != nil {
// 		return fmt.Errorf("could not register service Ping: %s", err)
// 	}
// 	log.Printf("starting HTTP/1.1 REST server on %s", address)
// 	http.ListenAndServe(address, mux)
// 	return nil
// }

// // unaryInterceptor - calls authenticateClient with current context
// func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	s, ok := info.Server.(*api.Server)
// 	if !ok {
// 		return nil, fmt.Errorf("unable to cast server")
// 	}
// 	clientID, err := authenticateClient(ctx, s)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ctx = context.WithValue(ctx, clientIDKey, clientID)
// 	return handler(ctx, req)
// }

func main() {
	grpcAddress := fmt.Sprintf("%s:%d", "localhost", 7777)
	// restAddress := fmt.Sprintf("%s:%d", "localhost", 7778)
	certFile := "cert/server.crt"
	keyFile := "cert/server.key"
	trustStore := []string{"cert/server.crt"}
	// // fire the gRPC server in a goroutine
	// go func() {
	err := startGRPCServer(grpcAddress, certFile, keyFile, trustStore)
	if err != nil {
		log.Fatalf("failed to start gRPC server: %s", err)
	}
	// }()
	// // // fire the REST server in a goroutine
	// // go func() {
	// // 	err := startRESTServer(restAddress, grpcAddress, certFile)
	// // 	if err != nil {
	// // 		log.Fatalf("failed to start gRPC server: %s", err)
	// // 	}
	// // }()
	// // // infinite loop
	// log.Printf("Entering infinite loop")
	select {}
}
