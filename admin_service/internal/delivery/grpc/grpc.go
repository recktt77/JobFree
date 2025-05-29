package grpc

import (
	"github.com/recktt77/JobFree/admin_service/internal/usecase"
	adminpb "github.com/recktt77/projectProto-definitions/gen/admin_service"
)

type AdminHandler struct {
	adminpb.UnimplementedAdminServiceServer
	UseCase usecase.AdminUseCase
}

func NewAdminHandler(uc usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		UseCase: uc,
	}
}
