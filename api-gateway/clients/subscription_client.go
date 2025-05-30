// clients/subscription_client.go
package clients

import (
	"log"
	"sync"

	pb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/subscription"
	"google.golang.org/grpc"
)

var (
	subscriptionClient pb.SubscriptionServiceClient
	subOnce            sync.Once
)

func GetSubscriptionClient() pb.SubscriptionServiceClient {
	subOnce.Do(func() {
		conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to SubscriptionService: %v", err)
		}
		subscriptionClient = pb.NewSubscriptionServiceClient(conn)
	})
	return subscriptionClient
}
