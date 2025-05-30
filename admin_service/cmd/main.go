package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/recktt77/JobFree/admin_service/internal/delivery"
	"github.com/recktt77/JobFree/admin_service/internal/events"
	"github.com/recktt77/JobFree/admin_service/internal/metrics"
	"github.com/recktt77/JobFree/admin_service/internal/repository"
	"github.com/recktt77/JobFree/admin_service/internal/usecase"
	adminpb "github.com/recktt77/projectProto-definitions/gen/admin_service"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	// MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Mongo connection error: %v", err)
	}
	db := client.Database("jobfree")
	log.Println("✅ Connected to MongoDB")

	// NATS (опционально)
	var nc *nats.Conn
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	nc, err = nats.Connect(natsURL)
	if err != nil {
		log.Printf("⚠️  NATS connection failed, continuing without it: %v", err)
		nc = nil
	} else {
		log.Println("✅ Connected to NATS")
	}

	// Prometheus metrics
	go func() {
		metrics.Register()
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	// Clean Architecture
	repo := repository.NewAdminRepository(db)
	publisher := events.NewEventPublisher(nc)
	uc := usecase.NewAdminUseCase(repo, publisher)
	handler := delivery.NewAdminHandler(uc)

	// gRPC Server
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	adminpb.RegisterAdminServiceServer(s, handler)

	log.Println("🚀 AdminService running on :8082")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}
