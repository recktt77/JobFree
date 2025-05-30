package clients

import (
	"log"
	"sync"

	pb "github.com/recktt77/projectProto-definitions/gen/admin_service"
	"google.golang.org/grpc"
)

var (
	adminClient pb.AdminServiceClient
	adminOnce   sync.Once
)

func GetAdminClient() pb.AdminServiceClient {
	adminOnce.Do(func() {
		conn, err := grpc.Dial("localhost:50058", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to AdminService: %v", err)
		}
		adminClient = pb.NewAdminServiceClient(conn)
	})
	return adminClient
}
