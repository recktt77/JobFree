package handler

import (
	"context"

	"github.com/recktt77/JobFree/matching_service/internal/usecase"

	matchingpb "github.com/recktt77/projectProto-definitions/gen/matching_service"
)

type MatchingHandler struct {
	matchingpb.UnimplementedMatchingServiceServer
	UseCase usecase.MatchingUseCase
}

func NewMatchingHandler(uc usecase.MatchingUseCase) *MatchingHandler {
	return &MatchingHandler{
		UseCase: uc,
	}
}

func (h *MatchingHandler) CreateBid(ctx context.Context, req *matchingpb.CreateBidRequest) (*matchingpb.CreateBidResponse, error) {
	// Заглушка для теста
	return &matchingpb.CreateBidResponse{BidId: "stub-id"}, nil
}

func (h *MatchingHandler) GetBidsForProject(ctx context.Context, req *matchingpb.GetBidsRequest) (*matchingpb.GetBidsResponse, error) {
	return &matchingpb.GetBidsResponse{}, nil
}

func (h *MatchingHandler) MatchFreelancers(ctx context.Context, req *matchingpb.MatchRequest) (*matchingpb.MatchResponse, error) {
	return &matchingpb.MatchResponse{}, nil
}
