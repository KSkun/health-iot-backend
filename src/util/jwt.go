package util

import (
	"github.com/KSkun/health-iot-backend/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func makeJWTClaims(id string) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Issuer:    "health-iot-backend",
		Subject:   id,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.C.JWTConfig.ExpireDuration)),
	}
}

func NewJWTToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, makeJWTClaims(id))
	return token.SignedString(config.C.JWTConfig.SecretBytes)
}

func ValidateJWTToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return config.C.JWTConfig.SecretBytes, nil
	})
	if err != nil {
		return "", err
	}
	return token.Claims.(jwt.RegisteredClaims).Subject, nil
}
