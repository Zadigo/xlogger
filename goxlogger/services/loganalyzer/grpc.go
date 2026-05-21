package loganalyzer

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	Port int32
}

func (s *GrpcServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatalf("Failed to listen on %d: %v", s.Port, err)
	}

	server := grpc.NewServer()
	return server.Serve(listener)
}

func NewGrpcServer(port int32) *GrpcServer {
	return &GrpcServer{Port: port}
}
