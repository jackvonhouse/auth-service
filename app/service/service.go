package service

import (
	"github.com/jackvonhouse/auth-service/app/repository"
	"github.com/jackvonhouse/auth-service/config"
	"github.com/jackvonhouse/auth-service/internal/service/auth"
	"github.com/jackvonhouse/auth-service/internal/service/jwt"
	"github.com/jackvonhouse/auth-service/pkg/log"
)

type Service struct {
	JWT  *jwt.ServiceJWT
	Auth *auth.ServiceAuth
}

func New(
	repository *repository.Repository,
	config *config.JWT,
	logger log.Logger,
) *Service {

	serviceLogger := logger.WithField("layer", "service")

	return &Service{
		JWT:  jwt.New(config, serviceLogger),
		Auth: auth.New(repository.Auth, config, logger),
	}
}
