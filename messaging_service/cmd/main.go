package main

import (
	"log"
	"net"

	"messaging_service/config"
	"messaging_service/email"
	"messaging_service/internal/broker"
	"messaging_service/internal/handler"
	"messaging_service/internal/repository"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/recktt77/projectProto-definitions/gen/messaging_service/genproto/messaging"
	"google.golang.org/grpc"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load")
	}

	// Get Mongo config from env
	mongoURI := config.GetEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := config.GetEnv("MONGO_DB_NAME", "messaging_db")

	// Init Mongo
	db, err := config.InitMongo(mongoURI, mongoDB)
	if err != nil {
		log.Fatalf("Mongo init error: %v", err)
	}
	repo := repository.NewMessageRepository(db)

	// SMTP config from env
	smtpSender := email.NewSmtpSender(
		config.GetEnv("SMTP_FROM", ""),
		config.GetEnv("SMTP_PASSWORD", ""),
		config.GetEnv("SMTP_HOST", "smtp.gmail.com"),
		config.GetEnv("SMTP_PORT", "587"),
	)

	// Init NATS connection
	natsURL := config.GetEnv("NATS_URL", nats.DefaultURL)
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("NATS connect error: %v", err)
	}

	// Subscribe to registration events
	natsSub := broker.NewNatsSubscriberCore(nc, repo, smtpSender)
	go func() {
		if err := natsSub.SubscribeToUserRegistered(); err != nil {
			log.Fatalf("NATS subscribe error: %v", err)
		}
	}()

	// gRPC config from env
	grpcPort := config.GetEnv("GRPC_PORT", ":50055")
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	srv := handler.NewMessagingHandler(repo)
	messaging.RegisterMessagingServiceServer(grpcServer, srv)

	log.Printf("MessagingService gRPC running on %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC serve error: %v", err)
	}
}
