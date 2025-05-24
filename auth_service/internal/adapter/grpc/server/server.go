package server

import (
    "github.com/recktt77/JobFree/config"
    "github.com/recktt77/projectProto-definitions/auth_service/genproto/auth"
    "github.com/recktt77/JobFree/internal/adapter/grpc/server/frontend"
    _"github.com/recktt77/JobFree/internal/usecase"
    "context"
    "fmt"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "log"
    "net"
)

type API struct {
    s           *grpc.Server
    cfg         config.GRPCServer
    addr        string
    userUsecase UserUsecase
}

func New(
    cfg config.GRPCServer,
    userUsecase UserUsecase,
) *API {
    return &API{
        cfg:         cfg,
        addr:        fmt.Sprintf("0.0.0.0:%d", cfg.Port),
        userUsecase: userUsecase,
    }
}

func (a *API) Run(ctx context.Context, errCh chan<- error) {
    go func() {
        log.Println("gRPC server listening at", a.addr)

        if err := a.run(ctx); err != nil {
            errCh <- fmt.Errorf("can't start grpc server: %w", err)
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

    auth.RegisterAuthServiceServer(a.s, frontend.NewUserHandler(a.userUsecase))

    reflection.Register(a.s)

    listener, err := net.Listen("tcp", a.addr)
    if err != nil {
        return fmt.Errorf("failed to create listener: %w", err)
    }

    return a.s.Serve(listener)
}
