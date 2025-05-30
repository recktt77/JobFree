package clients

import (
	"log"
	"sync"

	pb "github.com/recktt77/projectProto-definitions/gen/review_service/genproto/review"

	"google.golang.org/grpc"
)

var (
	reviewClient pb.ReviewServiceClient
	once         sync.Once
)

func GetReviewClient() pb.ReviewServiceClient {
	once.Do(func() {
		conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to ReviewService: %v", err)
		}
		reviewClient = pb.NewReviewServiceClient(conn)
	})
	return reviewClient
}
