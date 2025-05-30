package server

import (
	pb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/payment"
	"google.golang.org/grpc"
)

type Servers struct {
	Frontend FrontendServers
}

func NewServers(paymentHandler pb.PaymentServiceServer) *Servers {
	return &Servers{
		Frontend: FrontendServers{
			Payment: paymentHandler,
		},
	}
}



func (s *Servers) Register(server *grpc.Server) {
	pb.RegisterPaymentServiceServer(server, s.Frontend.Payment)
}
