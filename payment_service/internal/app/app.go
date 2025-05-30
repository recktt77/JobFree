package app

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/recktt77/JobFree/payment_service/internal/adapter/grpc/server"
	"github.com/recktt77/JobFree/payment_service/internal/adapter/grpc/server/frontend"
	"github.com/recktt77/JobFree/payment_service/internal/adapter/mongo"
	"github.com/recktt77/JobFree/payment_service/internal/adapter/nats/producer"
	"github.com/recktt77/JobFree/payment_service/internal/usecase"
	mongos "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Run() {
	// Mongo init
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongos.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	db := client.Database("jobfree")
	// NATS init
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}

	// Repository & Producer
	paymentRepo := mongo.NewPaymentRepository(db)
	paymentProducer := producer.NewPaymentProducer(nc)

	// Usecase
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, paymentProducer)

	// gRPC Handler
	paymentHandler := frontend.NewPaymentHandler(paymentUsecase)
	servers := server.NewServers(paymentHandler) // ✅ теперь всё совместимо

	if err := server.RunGRPCServer(servers, "50055"); err != nil {
		log.Fatal("gRPC server error:", err)
	}
}
