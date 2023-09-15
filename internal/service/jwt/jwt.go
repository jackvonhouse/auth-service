package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackvonhouse/auth-service/config"
	"github.com/jackvonhouse/auth-service/internal/dto"
	"github.com/jackvonhouse/auth-service/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTokenSize = 32
)

type ServiceJWT struct {
	logger    log.Logger
	secretKey string
	config    *config.JWT
}

func New(
	config *config.JWT,
	logger log.Logger,
) *ServiceJWT {

	return &ServiceJWT{
		logger:    logger.WithField("unit", "jwt"),
		secretKey: config.SecretKey,
		config:    config,
	}
}

func (s *ServiceJWT) CreateAccessToken(
	data *dto.AccessToken,
) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, AccessTokenClaim{
		GUID:           data.GUID,
		RefreshTokenId: data.RefreshTokenId,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(
					time.Duration(s.config.AccessToken.Exp) * time.Minute,
				),
			),
		},
	})

	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		s.logger.Warnf("cant sign access token: %s", err)

		return "", err
	}

	return signedToken, nil
}

func (s *ServiceJWT) CreateRefreshToken() (string, error) {
	buffer := make([]byte, refreshTokenSize)

	_, err := rand.Read(buffer)
	if err != nil {
		s.logger.Warnf("can't create refresh token: %s", err)

		return "", fmt.Errorf("ErrInternal: can't create refresh token")
	}

	b64 := base64.RawStdEncoding.EncodeToString(buffer)

	return b64, nil
}

func (s *ServiceJWT) getKey(t *jwt.Token) (interface{}, error) {
	return []byte(s.config.SecretKey), nil
}

func (s *ServiceJWT) ParseAccessToken(
	token string,
) (*dto.AccessToken, error) {

	claim := AccessTokenClaim{}

	_, err := jwt.ParseWithClaims(token, &claim, s.getKey)
	if err != nil {
		s.logger.Warnf("can't parse access token: %s", err)

		return nil, fmt.Errorf("access token has been modified")
	}

	return &dto.AccessToken{
		GUID:           claim.GUID,
		RefreshTokenId: claim.RefreshTokenId,
	}, nil
}

func (s *ServiceJWT) VerifyRefreshToken(
	token, hashedToken string,
) error {

	if err := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token)); err != nil {
		s.logger.Warnf("tokens not equals: %s", err)

		return fmt.Errorf("refresh token has been modified")
	}

	return nil
}
