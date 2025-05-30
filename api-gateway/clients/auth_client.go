package clients

import (
	"log"
	"sync"

	pb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/auth"

	"google.golang.org/grpc"
)

var (
	authClient pb.AuthServiceClient
	authOnce   sync.Once
)

func GetAuthClient() pb.AuthServiceClient {
	authOnce.Do(func() {
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to AuthService: %v", err)
		}
		authClient = pb.NewAuthServiceClient(conn)
	})
	return authClient
}
