package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackvonhouse/auth-service/config"
	"github.com/jackvonhouse/auth-service/internal/dto"
	"github.com/jackvonhouse/auth-service/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type RepositoryAuth interface {
	CreateRefreshToken(context.Context, *dto.RefreshToken) (string, error)
	GetRefreshToken(context.Context, string) (*dto.RefreshToken, error)
	DeleteRefreshToken(context.Context, string) error
}

type ServiceAuth struct {
	repository RepositoryAuth
	logger     log.Logger
	config     *config.JWT
}

func New(
	repository RepositoryAuth,
	config *config.JWT,
	logger log.Logger,
) *ServiceAuth {

	return &ServiceAuth{
		repository: repository,
		config:     config,
		logger:     logger.WithField("unit", "auth"),
	}
}

func (s *ServiceAuth) CreateRefreshToken(
	ctx context.Context,
	data *dto.RefreshToken,
) (string, error) {

	hashedToken, err := bcrypt.GenerateFromPassword(
		[]byte(data.Token),
		bcrypt.DefaultCost,
	)

	if err != nil {
		s.logger.Warnf("cant hash refresh token by bcrypt: %s", err)

		return "", fmt.Errorf("ErrInternal: cant hash refresh token by bcrypt")
	}

	data.Token = string(hashedToken)
	data.ExpireAt = time.Now().Add(
		time.Duration(s.config.RefreshToken.Exp) * time.Minute,
	)

	return s.repository.CreateRefreshToken(ctx, data)
}

func (s *ServiceAuth) GetRefreshToken(
	ctx context.Context,
	id string,
) (*dto.RefreshToken, error) {

	return s.repository.GetRefreshToken(ctx, id)
}

func (s *ServiceAuth) DeleteRefreshToken(
	ctx context.Context,
	id string,
) error {

	return s.repository.DeleteRefreshToken(ctx, id)
}
