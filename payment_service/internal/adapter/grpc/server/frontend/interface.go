package frontend

import (
	"context"

	pb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/payment"
)

type PaymentHandler interface {
	GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error)
	CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error)
	ListUserPayments(ctx context.Context, req *pb.ListUserPaymentsRequest) (*pb.ListUserPaymentsResponse, error)
}
