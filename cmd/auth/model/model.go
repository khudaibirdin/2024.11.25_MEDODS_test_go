package model

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthRefreshTokens struct {
	gorm.Model
	GUID         string `json:"guid"`
	RefreshToken string `json:"refresh_token"`
}

type GetTokensRequestStruct struct {
	GUID string `json:"guid"`
	IP   string `json:"ip"`
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	GUID string `json:"guid"`
	IP   string `json:"ip"`
}

type RefreshTokenRequestStruct struct {
	IP           string `json:"ip"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
