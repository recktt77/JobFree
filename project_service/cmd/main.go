// TODO: MIDDLEWARE, JWT TOKEN VALIDATION!!!!!!!!!!!!!!!!!!!!

package main

import (
	"log"
	"net"
	"project_service/config"
	"project_service/internal/cache"
	"project_service/internal/events"
	"project_service/internal/handler"
	"project_service/internal/repository"

	pb "github.com/recktt77/projectProto-definitions/gen/project_service/genproto/project"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	// Mongo
	mongo, err := repository.NewMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal("Mongo connection error:", err)
	}
	log.Println("Mongo connected")

	// Redis
	redis := cache.NewRedis(cfg.RedisAddr, cfg.RedisPwd)
	log.Println("Redis connected")

	// NATS
	nc, err := events.Connect(cfg.NatsURL)
	if err != nil {
		log.Fatal("NATS connection error:", err)
	}
	log.Println("NATS connected")

	// gRPC
	projectRepo := repository.NewProjectRepository(mongo.DB)
	publisher := events.NewPublisher(nc)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterProjectServiceServer(grpcServer, handler.New(projectRepo, publisher, redis))

	log.Println("gRPC server started on port", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
