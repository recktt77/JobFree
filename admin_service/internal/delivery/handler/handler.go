package handler

import (
	"context"

	"github.com/recktt77/JobFree/admin_service/internal/usecase"
	adminpb "github.com/recktt77/projectProto-definitions/gen/admin_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminHandler struct {
	adminpb.UnimplementedAdminServiceServer
	UseCase usecase.AdminUseCase
}

func NewAdminHandler(uc usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{UseCase: uc}
}

func (h *AdminHandler) BanUser(ctx context.Context, req *adminpb.BanUserRequest) (*adminpb.BanUserResponse, error) {
	if req.GetUserId() == "" || req.GetReason() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and reason are required")
	}

	err := h.UseCase.BanUser(ctx, req.GetUserId(), req.GetReason())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to ban user: %v", err)
	}

	return &adminpb.BanUserResponse{Success: true}, nil
}

func (h *AdminHandler) ApproveProject(ctx context.Context, req *adminpb.ModerateProjectRequest) (*adminpb.ModerateProjectResponse, error) {
	if req.GetProjectId() == "" || req.GetAction() == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id and action are required")
	}

	err := h.UseCase.ModerateProject(ctx, req.GetProjectId(), req.GetAction())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to moderate project: %v", err)
	}

	return &adminpb.ModerateProjectResponse{Success: true}, nil
}

func (h *AdminHandler) DeleteReview(ctx context.Context, req *adminpb.DeleteReviewRequest) (*adminpb.DeleteReviewResponse, error) {
	if req.GetReviewId() == "" || req.GetModeratorId() == "" {
		return nil, status.Error(codes.InvalidArgument, "review_id and moderator_id are required")
	}

	err := h.UseCase.DeleteReview(ctx, req.GetReviewId(), req.GetModeratorId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete review: %v", err)
	}

	return &adminpb.DeleteReviewResponse{Success: true}, nil
}

func (h *AdminHandler) GetPlatformStats(ctx context.Context, req *adminpb.GetStatsRequest) (*adminpb.GetStatsResponse, error) {
	stats, err := h.UseCase.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get stats: %v", err)
	}

	return &adminpb.GetStatsResponse{
		TotalUsers:     stats.TotalUsers,
		BannedUsers:    stats.BannedUsers,
		ActiveProjects: stats.ActiveProjects,
		TotalReviews:   stats.TotalReviews,
	}, nil
}
