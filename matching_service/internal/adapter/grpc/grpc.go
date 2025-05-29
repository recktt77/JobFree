package grpc

import (
	"github.com/recktt77/JobFree/matching_service/internal/adapter/grpc/handler"

	matchingpb "github.com/recktt77/projectProto-definitions/gen/matching_service"
)

func NewGrpcServer() matchingpb.MatchingServiceServer {
	return &handler.MatchingHandler{}
}
