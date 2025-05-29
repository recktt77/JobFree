package handler

import (
	"context"
	"fmt"

	"review_service/internal/broker"
	"review_service/internal/cache"
	"review_service/internal/model"
	"review_service/internal/repository"

	pb "github.com/recktt77/projectProto-definitions/gen/review_service/genproto/review"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewHandler struct {
	pb.UnimplementedReviewServiceServer
	repo   repository.ReviewRepository
	broker *broker.NatsBroker
	cache  *cache.ReviewCache
}

func NewReviewHandler(repo repository.ReviewRepository, broker *broker.NatsBroker, cache *cache.ReviewCache) *ReviewHandler {
	return &ReviewHandler{repo: repo, broker: broker, cache: cache}
}

func (h *ReviewHandler) LeaveReview(ctx context.Context, req *pb.LeaveReviewRequest) (*pb.LeaveReviewResponse, error) {
	// Простая валидация
	if req.ProjectId == "" || req.ReviewerId == "" || req.RevieweeId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id, reviewer_id and reviewee_id are required")
	}

	if req.Rating < 1 || req.Rating > 5 {
		return nil, status.Error(codes.InvalidArgument, "rating must be between 1 and 5")
	}

	// Попытка преобразования ID
	projectID, err := primitive.ObjectIDFromHex(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid project_id")
	}
	reviewerID, err := primitive.ObjectIDFromHex(req.ReviewerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reviewer_id")
	}
	revieweeID, err := primitive.ObjectIDFromHex(req.RevieweeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reviewee_id")
	}

	// Проверка на одинаковые reviewer и reviewee
	if reviewerID == revieweeID {
		return nil, status.Error(codes.InvalidArgument, "reviewer_id and reviewee_id must be different")
	}

	review := &model.Review{
		ProjectID:  projectID,
		ReviewerID: reviewerID,
		RevieweeID: revieweeID,
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	_, err = h.repo.Create(ctx, review)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create review: %v", err)
	}

	h.broker.PublishReviewCreated(review)
	h.cache.Invalidate(ctx, req.ProjectId, req.RevieweeId)

	return &pb.LeaveReviewResponse{Message: "Review submitted"}, nil
}

func (h *ReviewHandler) GetReviews(ctx context.Context, req *pb.GetReviewsRequest) (*pb.GetReviewsResponse, error) {
	var reviews []model.Review

	if req.ProjectId == "" && req.RevieweeId == "" {
		return nil, status.Error(codes.InvalidArgument, "either project_id or reviewee_id must be provided")
	}

	// Попытка получить из кэша
	if req.ProjectId != "" {
		projectID, err := primitive.ObjectIDFromHex(req.ProjectId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid project_id")
		}

		reviews, err = h.cache.GetProjectReviews(ctx, req.ProjectId)
		if err == nil && len(reviews) > 0 {
			fmt.Println("[CACHE HIT] project_id:", req.ProjectId)
			goto convert
		}

		reviews, err = h.repo.GetByProjectID(ctx, projectID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get from Mongo: %v", err)
		}
		_ = h.cache.SetProjectReviews(ctx, req.ProjectId, reviews)

	} else if req.RevieweeId != "" {
		revieweeID, err := primitive.ObjectIDFromHex(req.RevieweeId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid reviewee_id")
		}

		reviews, err = h.cache.GetUserReviews(ctx, req.RevieweeId)
		if err == nil && len(reviews) > 0 {
			fmt.Println("[CACHE HIT] reviewee_id:", req.RevieweeId)
			goto convert
		}

		reviews, err = h.repo.GetByRevieweeID(ctx, revieweeID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get from Mongo: %v", err)
		}
		_ = h.cache.SetUserReviews(ctx, req.RevieweeId, reviews)
	}

convert:
	var resp []*pb.Review
	for _, r := range reviews {
		resp = append(resp, &pb.Review{
			Id:         r.ID.Hex(),
			ProjectId:  r.ProjectID.Hex(),
			ReviewerId: r.ReviewerID.Hex(),
			RevieweeId: r.RevieweeID.Hex(),
			Rating:     r.Rating,
			Comment:    r.Comment,
		})
	}

	return &pb.GetReviewsResponse{Reviews: resp}, nil
}

func (h *ReviewHandler) ModerateReview(ctx context.Context, req *pb.ModerateReviewRequest) (*pb.ModerateReviewResponse, error) {
	if req.ReviewId == "" {
		return nil, status.Error(codes.InvalidArgument, "review_id is required")
	}
	if req.Action != "delete" && req.Action != "hide" && req.Action != "approve" {
		return nil, status.Error(codes.InvalidArgument, "invalid action, must be 'delete', 'hide', or 'approve'")
	}

	reviewID, err := primitive.ObjectIDFromHex(req.ReviewId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid review_id")
	}

	// Получаем отзыв перед удалением/скрытием
	review, err := h.repo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "review not found: %v", err)
	}

	// Применяем действие
	err = h.repo.Moderate(ctx, reviewID, req.Action)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to apply action: %v", err)
	}

	// Публикация и инвалидирование
	h.broker.PublishReviewModerated(reviewID.Hex(), req.Action)
	h.cache.Invalidate(ctx, review.ProjectID.Hex(), review.RevieweeID.Hex())

	return &pb.ModerateReviewResponse{Message: "Action applied"}, nil
}
