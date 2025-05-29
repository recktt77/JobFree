package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"review_service/config"
	"review_service/internal/broker"
	"review_service/internal/cache"
	"review_service/internal/handler"
	"review_service/internal/repository"

	pb "github.com/recktt77/projectProto-definitions/gen/review_service/genproto/review"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// MongoDB connection
	db, err := config.SafeConnectMongo()
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	repo := repository.NewReviewRepo(db)

	// Redis connection
	cache, err := cache.NewReviewCache()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	// NATS connection
	broker, err := broker.NewNatsBroker()
	if err != nil {
		log.Fatalf("NATS connection failed: %v", err)
	}

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	srv := handler.NewReviewHandler(repo, broker, cache)
	pb.RegisterReviewServiceServer(grpcServer, srv)

	fmt.Println("ReviewService gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
