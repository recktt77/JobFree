package main

import (
	"log"
	"net"

	"context"
	"os"
	"time"

	"github.com/recktt77/JobFree/admin_service/internal/adapter/grpc"
	"github.com/recktt77/JobFree/admin_service/internal/repository"
	"github.com/recktt77/JobFree/admin_service/internal/usecase"
	adminpb "github.com/recktt77/projectProto-definitions/gen/admin_service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("jobfree_admin")

	repo := repository.NewAdminRepository(db)
	uc := usecase.NewAdminUseCase(repo)
	handler := grpc.NewAdminHandler(uc)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	adminpb.RegisterAdminServiceServer(s, handler)

	log.Println("AdminService is running on port :8082")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
