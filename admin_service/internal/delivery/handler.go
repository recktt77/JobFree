package delivery

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
	err := h.UseCase.BanUser(ctx, req.GetUserId(), req.GetReason())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "BanUser error: %v", err)
	}
	return &adminpb.BanUserResponse{Success: true}, nil
}

func (h *AdminHandler) DeleteReview(ctx context.Context, req *adminpb.DeleteReviewRequest) (*adminpb.DeleteReviewResponse, error) {
	err := h.UseCase.DeleteReview(ctx, req.GetReviewId(), req.GetModeratorId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteReview error: %v", err)
	}
	return &adminpb.DeleteReviewResponse{Success: true}, nil
}

func (h *AdminHandler) ModerateProject(ctx context.Context, req *adminpb.ModerateProjectRequest) (*adminpb.ModerateProjectResponse, error) {
	err := h.UseCase.ModerateProject(ctx, req.GetProjectId(), req.GetAction())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModerateProject error: %v", err)
	}
	return &adminpb.ModerateProjectResponse{Success: true}, nil
}

func (h *AdminHandler) GetPlatformStats(ctx context.Context, req *adminpb.GetStatsRequest) (*adminpb.GetStatsResponse, error) {
	stats, err := h.UseCase.GetPlatformStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetStats error: %v", err)
	}
	return &adminpb.GetStatsResponse{
		TotalUsers:     stats.TotalUsers,
		BannedUsers:    stats.BannedUsers,
		ActiveProjects: stats.ActiveProjects,
		TotalReviews:   stats.TotalReviews,
	}, nil
}
