package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/recktt77/JobFree/matching_service/internal/adapter/grpc/handler"

	"github.com/recktt77/JobFree/matching_service/internal/repository"
	"github.com/recktt77/JobFree/matching_service/internal/usecase"

	"github.com/recktt77/JobFree/matching_service/internal/config"

	"github.com/joho/godotenv"

	matchingpb "github.com/recktt77/projectProto-definitions/gen/matching_service"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default env vars")
	}

	// Read env variables
	mongoURI := os.Getenv("MONGO_DB_URI")
	mongoDBName := os.Getenv("MONGO_DB")

	// Connect to MongoDB
	ctx := context.Background()
	mongoCfg := config.Config{
		URI:      mongoURI,
		Database: mongoDBName,
	}

	db, err := config.NewDB(ctx, mongoCfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create repository and usecase
	var mongoDB *mongo.Database = db.Conn
	bidRepo := repository.NewBidRepository(mongoDB)
	uc := usecase.NewMatchingUseCase(bidRepo)
	grpcHandler := handler.NewMatchingHandler(uc)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	matchingpb.RegisterMatchingServiceServer(server, grpcHandler)

	fmt.Println("MatchingService gRPC server is running on port 8081")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
