package usecase

import (
	"github.com/jackvonhouse/auth-service/app/service"
	"github.com/jackvonhouse/auth-service/internal/usecase/auth"
	"github.com/jackvonhouse/auth-service/pkg/log"
)

type UseCase struct {
	Auth *auth.UseCaseAuth
}

func New(
	service *service.Service,
	logger log.Logger,
) *UseCase {

	useCaseLogger := logger.WithField("layer", "usecase")

	return &UseCase{
		Auth: auth.New(
			service.JWT,
			service.Auth,
			useCaseLogger,
		),
	}
}
