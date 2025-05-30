package server

import (
	pb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/payment"
)

type FrontendServers struct {
	Payment pb.PaymentServiceServer
}
