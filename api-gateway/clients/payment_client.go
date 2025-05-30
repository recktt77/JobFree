// clients/payment_client.go
package clients

import (
	"log"
	"sync"

	pb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/payment"
	"google.golang.org/grpc"
)

var (
	paymentClient pb.PaymentServiceClient
	paymentOnce   sync.Once
)

func GetPaymentClient() pb.PaymentServiceClient {
	paymentOnce.Do(func() {
		conn, err := grpc.Dial("localhost:50057", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to PaymentService: %v", err)
		}
		paymentClient = pb.NewPaymentServiceClient(conn)
	})
	return paymentClient
}
