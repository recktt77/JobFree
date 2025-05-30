package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/recktt77/JobFree/subscription_service/config"
	"github.com/recktt77/JobFree/subscription_service/internal/usecase"
	"github.com/recktt77/JobFree/subscription_service/internal/adapter/grpc/server/frontend"
	subsrv "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/subscription"
)

type API struct {
	s        *grpc.Server
	cfg      config.GRPCServer
	addr     string
	usecase  usecase.SubscriptionUsecase
}

func New(cfg config.GRPCServer, uc usecase.SubscriptionUsecase) *API {
	return &API{
		cfg:     cfg,
		addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		usecase: uc,
	}
}

func (a *API) Run(ctx context.Context, errCh chan<- error) {
	go func() {
		log.Println("Subscription gRPC server listening on", a.addr)

		if err := a.run(ctx); err != nil {
			errCh <- fmt.Errorf("cannot start subscription grpc server: %w", err)
			return
		}
	}()
}

func (a *API) Stop(ctx context.Context) error {
	if a.s == nil {
		return nil
	}
	stopped := make(chan struct{})
	go func() {
		a.s.GracefulStop()
		close(stopped)
	}()
	select {
	case <-ctx.Done():
		a.s.Stop()
	case <-stopped:
	}
	return nil
}

func (a *API) run(ctx context.Context) error {
	a.s = grpc.NewServer()

	// регистрация gRPC-сервиса
	subsrv.RegisterSubscriptionServiceServer(a.s, frontend.NewSubscription(a.usecase))

	reflection.Register(a.s)

	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	if err := a.s.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve grpc: %w", err)
	}
	return nil
}
