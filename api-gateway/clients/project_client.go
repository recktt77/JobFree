package clients

import (
	"log"
	"sync"

	pb "github.com/recktt77/projectProto-definitions/gen/project_service/genproto/project"
	"google.golang.org/grpc"
)

var (
	projectClient pb.ProjectServiceClient
	projectOnce   sync.Once
)

func GetProjectClient() pb.ProjectServiceClient {
	projectOnce.Do(func() {
		conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to ProjectService: %v", err)
		}
		projectClient = pb.NewProjectServiceClient(conn)
	})
	return projectClient
}
