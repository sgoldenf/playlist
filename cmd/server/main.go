package main

import (
	"flag"
	"fmt"
	ps "github.com/sgoldenf/playlist/api"
	"github.com/sgoldenf/playlist/internal/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

var port = flag.Int("port", 50051, "gRPC server port")

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service, errService := server.NewService()
	if errService != nil {
		log.Fatalf("Failed to serve Database: %v", errService)
	}
	ps.RegisterPlaylistServiceServer(s, service)
	log.Fatalf("Failed to serve: %v", s.Serve(lis))
}
