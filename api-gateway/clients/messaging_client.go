package clients

import (
	"log"
	"sync"

	pb "github.com/recktt77/projectProto-definitions/gen/messaging_service/genproto/messaging"

	"google.golang.org/grpc"
)

var (
	messagingClient pb.MessagingServiceClient
	messagingOnce   sync.Once
)

func GetMessagingClient() pb.MessagingServiceClient {
	messagingOnce.Do(func() {
		conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to MessagingService: %v", err)
		}
		messagingClient = pb.NewMessagingServiceClient(conn)
	})
	return messagingClient
}
