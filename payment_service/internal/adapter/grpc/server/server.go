package server

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunGRPCServer(servers *Servers, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	servers.Register(grpcServer)

	log.Printf("gRPC server started on port %s\n", port)

	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
