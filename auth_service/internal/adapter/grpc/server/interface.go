package server

import "github.com/recktt77/JobFree/internal/adapter/grpc/server/frontend"

type UserUsecase interface {
	frontend.UserUsecase
}