package clients

import (
	"log"
	"sync"

	pb "github.com/recktt77/projectProto-definitions/gen/matching_service/recktt77/projectProto-definitions/matching_service"

	"google.golang.org/grpc"
)

var (
	matchingClient pb.MatchingServiceClient
	matchingOnce   sync.Once
)

func GetMatchingClient() pb.MatchingServiceClient {
	matchingOnce.Do(func() {
		conn, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to MatchingService: %v", err)
		}
		matchingClient = pb.NewMatchingServiceClient(conn)
	})
	return matchingClient
}
