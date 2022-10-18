package util

import (
	"github.com/KSkun/health-iot-backend/config"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

func makeJWTClaims(id string, issueTime time.Time, expireTime time.Time) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Issuer:    "health-iot-backend",
		Subject:   id,
		IssuedAt:  jwt.NewNumericDate(issueTime),
		ExpiresAt: jwt.NewNumericDate(expireTime),
	}
}

func NewJWTToken(id string) (string, time.Time, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(config.C.JWTConfig.ExpireDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, makeJWTClaims(id, nowTime, expireTime))
	tokenStr, err := token.SignedString(config.C.JWTConfig.SecretBytes)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenStr, expireTime, nil
}

func ValidateJWTToken(tokenStr string) (string, error) {
	options := []jwt.ParserOption{}
	if os.Getenv("JWT_NO_VALIDATE") == "1" {
		options = append(options, jwt.WithoutClaimsValidation())
	}
	claims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return config.C.JWTConfig.SecretBytes, nil
	}, options...)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}
