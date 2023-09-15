package jwt

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaim struct {
	GUID           string `json:"guid"`
	RefreshTokenId string `json:"id"`

	jwt.RegisteredClaims
}
