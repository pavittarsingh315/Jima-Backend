package utils

import (
	"NeraJima/configs"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenAuthTokens(user_id string) (access, refresh string) {
	accessSecret, refreshSecret := configs.EnvTokenSecrets()

	type Claims struct {
		Type   string `json:"token_type"`
		UserId string `json:"userId"`
		jwt.RegisteredClaims
	}

	accessExpTime := time.Now().Add(time.Hour * (24 * 30))       // 30 days
	refreshExpTime := time.Now().Add(time.Hour * (24 * 365 * 2)) // 2 Years

	accessClaims := Claims{
		"access",
		user_id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshClaims := Claims{
		"refresh",
		user_id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessSigningKey := []byte(accessSecret)
	refreshSigningKey := []byte(refreshSecret)

	accessSigned, _ := accessToken.SignedString(accessSigningKey)
	refreshSigned, _ := refreshToken.SignedString(refreshSigningKey)

	return accessSigned, refreshSigned
}
