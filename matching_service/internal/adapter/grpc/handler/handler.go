package handler

import (
	"context"
	"strings"

	"github.com/recktt77/JobFree/matching_service/internal/model"
	"github.com/recktt77/JobFree/matching_service/internal/usecase"

	matchingpb "github.com/recktt77/projectProto-definitions/gen/matching_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// CreateBid handles bid creation with validation
func (h *MatchingHandler) CreateBid(ctx context.Context, req *matchingpb.CreateBidRequest) (*matchingpb.CreateBidResponse, error) {
	// 🔒 Валидация
	if isEmpty(req.GetProjectId(), req.GetFreelancerId()) {
		return nil, status.Error(codes.InvalidArgument, "project_id and freelancer_id are required")
	}
	if len(strings.TrimSpace(req.GetCoverLetter())) < 10 {
		return nil, status.Error(codes.InvalidArgument, "cover_letter must be at least 10 characters")
	}

	// 🧠 Бизнес-логика
	bid := &model.Bid{
		BidID:        generateID(),
		ProjectID:    req.GetProjectId(),
		FreelancerID: req.GetFreelancerId(),
		CoverLetter:  req.GetCoverLetter(),
	}
	id, err := h.UseCase.CreateBid(ctx, bid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create bid: %v", err)
	}
	return &matchingpb.CreateBidResponse{BidId: id}, nil
}

// GetBidsForProject handles getting all bids for a project with validation
func (h *MatchingHandler) GetBidsForProject(ctx context.Context, req *matchingpb.GetBidsRequest) (*matchingpb.GetBidsResponse, error) {
	if strings.TrimSpace(req.GetProjectId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id is required")
	}

	bids, err := h.UseCase.GetBidsForProject(ctx, req.GetProjectId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get bids: %v", err)
	}

	var pbBids []*matchingpb.Bid
	for _, b := range bids {
		pbBids = append(pbBids, &matchingpb.Bid{
			BidId:        b.BidID,
			FreelancerId: b.FreelancerID,
			CoverLetter:  b.CoverLetter,
			Timestamp:    b.Timestamp,
		})
	}

	return &matchingpb.GetBidsResponse{Bids: pbBids}, nil
}

// MatchFreelancers handles freelancer matching (basic skill logic)
func (h *MatchingHandler) MatchFreelancers(ctx context.Context, req *matchingpb.MatchRequest) (*matchingpb.MatchResponse, error) {
	if strings.TrimSpace(req.GetProjectId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id is required")
	}

	requiredSkills := []string{"Go", "MongoDB"}

	freelancers, err := h.UseCase.MatchFreelancers(ctx, requiredSkills)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to match freelancers: %v", err)
	}

	var resp []*matchingpb.Freelancer
	for _, f := range freelancers {
		resp = append(resp, &matchingpb.Freelancer{
			FreelancerId: f.FreelancerID,
			Name:         f.Name,
			Skills:       f.Skills,
		})
	}

	return &matchingpb.MatchResponse{Freelancers: resp}, nil
}

// generateID временно
func generateID() string {
	return "bid_" + model.GenerateUUID()
}

// isEmpty checks for blank strings
func isEmpty(strs ...string) bool {
	for _, s := range strs {
		if strings.TrimSpace(s) == "" {
			return true
		}
	}
	return false
}
