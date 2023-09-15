package dto

import "time"

type Auth struct {
	GUID string `json:"guid"`
}

type RefreshToken struct {
	Token    string    `bson:"token"`
	ExpireAt time.Time `bson:"expire_at"`
}

type AccessToken struct {
	GUID           string
	RefreshTokenId string
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
