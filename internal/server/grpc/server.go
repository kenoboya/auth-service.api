package grpc_server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	srv *grpc.Server
}

func NewServer() *server {
	return &server{
		srv: grpc.NewServer(),
	}
}

func (s *server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on port %s : %v", addr, err)
		return err
	}

	if err := s.srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %s : %v", addr, err)
	}

	return nil
}
