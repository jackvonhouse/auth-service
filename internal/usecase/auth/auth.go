package auth

import (
	"context"

	"github.com/jackvonhouse/auth-service/internal/dto"
	"github.com/jackvonhouse/auth-service/pkg/log"
)

type serviceJWT interface {
	CreateAccessToken(*dto.AccessToken) (string, error)
	CreateRefreshToken() (string, error)

	ParseAccessToken(string) (*dto.AccessToken, error)

	VerifyRefreshToken(string, string) error
}

type serviceAuth interface {
	CreateRefreshToken(context.Context, *dto.RefreshToken) (string, error)
	GetRefreshToken(context.Context, string) (*dto.RefreshToken, error)
	DeleteRefreshToken(context.Context, string) error
}

type UseCaseAuth struct {
	jwt  serviceJWT
	auth serviceAuth

	logger log.Logger
}

func New(
	jwt serviceJWT,
	auth serviceAuth,
	logger log.Logger,
) *UseCaseAuth {

	return &UseCaseAuth{
		jwt:    jwt,
		auth:   auth,
		logger: logger.WithField("unit", "auth"),
	}
}

func (u *UseCaseAuth) Auth(
	ctx context.Context,
	data *dto.Auth,
) (*dto.TokenPair, error) {

	refreshToken, err := u.jwt.CreateRefreshToken()
	if err != nil {
		return nil, err
	}

	id, err := u.auth.CreateRefreshToken(ctx, &dto.RefreshToken{
		Token: refreshToken,
	})

	if err != nil {
		return nil, err
	}

	accessToken, err := u.jwt.CreateAccessToken(&dto.AccessToken{
		GUID:           data.GUID,
		RefreshTokenId: id,
	})

	if err != nil {
		return nil, err
	}

	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UseCaseAuth) Refresh(
	ctx context.Context,
	data *dto.TokenPair,
) (*dto.TokenPair, error) {

	accessToken, err := u.jwt.ParseAccessToken(data.AccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.auth.GetRefreshToken(ctx, accessToken.RefreshTokenId)
	if err != nil {
		return nil, err
	}

	if err := u.jwt.VerifyRefreshToken(data.RefreshToken, refreshToken.Token); err != nil {
		return nil, err
	}

	if err := u.auth.DeleteRefreshToken(ctx, accessToken.RefreshTokenId); err != nil {
		return nil, err
	}

	return u.Auth(ctx, &dto.Auth{
		GUID: accessToken.GUID,
	})
}
