package app

import (
	"context"
	"fmt"
	_ "fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/recktt77/JobFree/config"
	"github.com/recktt77/JobFree/internal/adapter/grpc/server"
	_ "github.com/recktt77/JobFree/internal/adapter/inmemory"
	"github.com/recktt77/JobFree/internal/adapter/mongo"
	"github.com/recktt77/JobFree/internal/adapter/mongo/dao"
	"github.com/recktt77/JobFree/internal/adapter/nats/producer"
	rediscache "github.com/recktt77/JobFree/internal/adapter/redis"
	"github.com/recktt77/JobFree/internal/usecase"
	mongopkg "github.com/recktt77/JobFree/pkg/mongo"
	natspkg "github.com/recktt77/JobFree/pkg/nats"
	redispkg "github.com/recktt77/JobFree/pkg/redis"
)

const shutdownTimeout = 30 * time.Second

type App struct {
    grpcServer  *server.API
    redisClient *rediscache.RedisCache
}

func Run(ctx context.Context, cfg config.Config) {
    log.Println("Starting user-service...")

    // Mongo
    mongoClient, err := mongopkg.NewDB(ctx, cfg.Mongo)
    if err != nil {
        log.Fatalf("mongo connection failed: %v", err)
    }
    userDAO := dao.NewUserDAO(mongoClient.Conn)
    userRepo := mongo.NewUserRepository(userDAO)

    // Redis
    redisClient, err := redispkg.NewProduct(ctx, redispkg.Config(cfg.Redis))
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}

	redisCache := rediscache.NewRedisCache(redisClient.Unwrap(), cfg.Cache.UserTTL)




    // InMemory fallback
    // memoryCache := inmemory.NewInMemoryCache(cfg.Cache.UserTTL)

    // // Example: initialize in-memory cache with all users (if you have a GetAll method)
    // users, err := userRepo.GetAll(ctx)
    // if err == nil {
    //     memoryCache.SetMany(users)
    //     log.Println("in-memory user cache initialized from DB")
    // } else {
    //     log.Println("failed to init in-memory cache:", err)
    // }

    // NATS
    natsConn, err := natspkg.Connect(cfg.NATSUrl)
    if err != nil {
        log.Fatalf("NATS connection failed: %v", err)
    }
    eventPublisher := producer.NewUserEventPublisher(natsConn)

    // UseCase
    userUC := usecase.NewUserUseCase(userRepo, redisCache, eventPublisher)

    // gRPC Server
    grpcServer := server.New(cfg.Server.GRPCServer, userUC)

    // Run server
    errCh := make(chan error, 1)
    grpcServer.Run(ctx, errCh)

    shutdownCh := make(chan struct{})
    go func() {
        select {
        case err := <-errCh:
            log.Fatalf("server error: %v", err)
        case <-ctx.Done():
            log.Println("Shutdown initiated...")
            shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
            defer cancel()
            _ = grpcServer.Stop(shutdownCtx)
            shutdownCh <- struct{}{}
        }
    }()

    <-shutdownCh
    log.Println("Shutdown complete.")
}
func (a *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()
	go a.grpcServer.Run(ctx, errCh)
	log.Println(fmt.Sprintf("service %v started (gRPC)", "auth-service"))

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
	if err := a.redisClient.Close(); err != nil {
		log.Println("failed to close redis:", err)
	}

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