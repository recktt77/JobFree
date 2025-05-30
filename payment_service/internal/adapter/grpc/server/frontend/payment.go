package frontend

import (
	"context"

	"github.com/recktt77/JobFree/payment_service/internal/adapter/grpc/server/frontend/dto"
	"github.com/recktt77/JobFree/payment_service/internal/model"
	"github.com/recktt77/JobFree/payment_service/internal/usecase"
	pb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/payment"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type paymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	usecase *usecase.PaymentUsecase
}

func NewPaymentHandler(u *usecase.PaymentUsecase) *paymentHandler {
	return &paymentHandler{
		usecase: u,
	}
}



func (h *paymentHandler) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	payment, err := h.usecase.GetPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return &pb.GetPaymentResponse{}, nil
	}
	return &pb.GetPaymentResponse{
		Payment: dto.ToProto(payment),
	}, nil
}

func (h *paymentHandler) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	p := model.Payment{
		UserID:         ObjectIDFromHexOrNil(req.UserId),
		SubscriptionID: ObjectIDFromHexOrNil(req.SubscriptionId),
		Status:         "pending",
	}

	err := h.usecase.CreatePayment(ctx, p)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		PaymentId: p.ID.Hex(),
	}, nil

}


func (h *paymentHandler) ListUserPayments(ctx context.Context, req *pb.ListUserPaymentsRequest) (*pb.ListUserPaymentsResponse, error) {
	payments, err := h.usecase.ListUserPayments(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var pbPayments []*pb.Payment
	for _, p := range payments {
		pbPayments = append(pbPayments, dto.ToProto(&p))
	}

	return &pb.ListUserPaymentsResponse{
		Payments: pbPayments,
	}, nil
}

func ObjectIDFromHexOrNil(hex string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return primitive.NilObjectID
	}
	return id
}
