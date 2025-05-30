package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/recktt77/JobFree/subscription_service/config"
	"github.com/recktt77/JobFree/subscription_service/internal/adapter/grpc/server"
	inmemorycache "github.com/recktt77/JobFree/subscription_service/internal/adapter/inmemory"
	mongorepo "github.com/recktt77/JobFree/subscription_service/internal/adapter/mongo"
	_"github.com/recktt77/JobFree/subscription_service/internal/adapter/mongo/dao"
	"github.com/recktt77/JobFree/subscription_service/internal/adapter/nats/producer"
	rediscache "github.com/recktt77/JobFree/subscription_service/internal/adapter/redis"
	"github.com/recktt77/JobFree/subscription_service/internal/usecase"
	mongocon "github.com/recktt77/JobFree/subscription_service/pkg/mongo"
	redispkg "github.com/recktt77/JobFree/subscription_service/pkg/redis"
)

const (
	serviceName     = "subscription-service"
	shutdownTimeout = 30 * time.Second
)

type App struct {
	grpcServer    *server.API
	redisClient   *redispkg.Subscription
	subCache      *rediscache.Subscription
	inMemoryCache *inmemorycache.Subscription
}

func New(ctx context.Context, cfg *config.Config, natsConn *nats.Conn) (*App, error) {
	log.Println(fmt.Sprintf("starting %s", serviceName))

	// MongoDB
	mongoDB, err := mongocon.NewDB(ctx, cfg.Mongo)
	if err != nil {
		return nil, fmt.Errorf("mongo init failed: %w", err)
	}

	// Redis
	redisClient, err := redispkg.NewProduct(ctx, redispkg.Config(cfg.Redis))
	if err != nil {
		return nil, fmt.Errorf("redis init failed: %w", err)
	}

	// DAO / Repository
	subRepo := mongorepo.NewRepository(mongoDB.Conn)

	// Caches
	inMemory := inmemorycache.NewSubscription()
	redisCache := rediscache.NewSubscription(redisClient, cfg.Cache.SubscriptionTTL)

	// NATS producer
	prod := producer.NewProducer(natsConn)

	// Usecase
	subUC := usecase.NewSubscriptionUsecase(subRepo, prod)

	// gRPC server
	grpcSrv := server.New(cfg.Server.GRPCServer, subUC)

	return &App{
		grpcServer:    grpcSrv,
		redisClient:   redisClient,
		subCache:      redisCache,
		inMemoryCache: inMemory,
	}, nil
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()

	go a.grpcServer.Run(ctx, errCh)
	log.Println(fmt.Sprintf("service %s started (gRPC)", serviceName))

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("gRPC server error: %w", err)
	case sig := <-shutdownCh:
		log.Println(fmt.Sprintf("received signal: %v. Shutting down...", sig))
		return a.Close()
	}
}

func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.grpcServer.Stop(ctx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("shutdown timed out after %v", shutdownTimeout)
	}

	log.Println("graceful shutdown completed successfully")
	return nil
}
