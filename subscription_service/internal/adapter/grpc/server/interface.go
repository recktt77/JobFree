package server

import "github.com/recktt77/JobFree/subscription_service/internal/adapter/grpc/server/frontend"

type SubscriptionUsecase interface {
	frontend.Handler
}