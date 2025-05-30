package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/recktt77/JobFree/matching_service/internal/adapter/grpc/handler"
	"github.com/recktt77/JobFree/matching_service/internal/cache"
	"github.com/recktt77/JobFree/matching_service/internal/config"
	"github.com/recktt77/JobFree/matching_service/internal/events"
	"github.com/recktt77/JobFree/matching_service/internal/repository"
	"github.com/recktt77/JobFree/matching_service/internal/seed"
	"github.com/recktt77/JobFree/matching_service/internal/usecase"
	matchingpb "github.com/recktt77/projectProto-definitions/gen/matching_service"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	mongoCfg := config.Config{
		URI:      os.Getenv("MONGO_DB_URI"),
		Database: os.Getenv("MONGO_DB"),
	}
	db, err := config.NewDB(ctx, mongoCfg)
	if err != nil {
		log.Fatalf("❌ failed to connect to MongoDB: %v", err)
	}
	var mongoDB *mongo.Database = db.Conn

	redisCache := cache.NewRedisCache("localhost:6379")
	seed.Run(redisCache)

	publisher := events.NewPublisher("nats://localhost:4222")

	bidRepo := repository.NewBidRepository(mongoDB)
	uc := usecase.NewMatchingUseCase(bidRepo, redisCache, publisher)

	grpcHandler := handler.NewMatchingHandler(uc)
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("❌ failed to listen: %v", err)
	}

	server := grpc.NewServer()
	matchingpb.RegisterMatchingServiceServer(server, grpcHandler)

	fmt.Println("🚀 MatchingService is running on :50054")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("❌ failed to serve: %v", err)
	}
}
