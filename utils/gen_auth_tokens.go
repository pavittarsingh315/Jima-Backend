package utils

import (
	"NeraJima/configs"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	Type   string `json:"token_type"`
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func GenAuthTokens(user_id string) (access, refresh string) {
	accessSecret, refreshSecret := configs.EnvTokenSecrets()

	accessExpTime := time.Now().Add(time.Hour * (24 * 30))       // 30 days
	refreshExpTime := time.Now().Add(time.Hour * (24 * 365 * 2)) // 2 Years

	accessClaims := claims{
		"access",
		user_id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshClaims := claims{
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

func VerifyAccessToken(token string) (string, claims, error) {
	accessSecret, _ := configs.EnvTokenSecrets()
	var tokenBody claims

	_, err := jwt.ParseWithClaims(token, &tokenBody, func(t *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})

	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			accessExpTime := time.Now().Add(time.Hour * (24 * 30)) // 30 days
			tokenBody.IssuedAt = jwt.NewNumericDate(time.Now())
			tokenBody.ExpiresAt = jwt.NewNumericDate(accessExpTime)
			newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenBody)
			accessSigningKey := []byte(accessSecret)
			newToken, _ := newAccessToken.SignedString(accessSigningKey)
			return newToken, tokenBody, nil
		} else {
			return "", claims{}, err
		}
	} else {
		return token, tokenBody, nil
	}
}

func VerifyRefreshToken(token string) (string, claims, error) {
	_, refreshSecret := configs.EnvTokenSecrets()
	var tokenBody claims

	_, err := jwt.ParseWithClaims(token, &tokenBody, func(t *jwt.Token) (interface{}, error) {
		return []byte(refreshSecret), nil
	})

	if err != nil {
		return "", claims{}, err
	} else {
		return token, tokenBody, nil
	}
}
