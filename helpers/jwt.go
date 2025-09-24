package helpers

import (
	"context"
	"strings"
	"time"
	errConstants "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/helpers/configs"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimToken struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func GenerateToken(ctx context.Context, user *dto.UserResponse, expirationTime int64) (string, error) {
	claimToken := &ClaimToken{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	cfg := configs.Get()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)

	tokenString, err := token.SignedString([]byte(cfg.JwtConfig.JwtSecretKey))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")
	if len(arrayToken) == 2 {
		return arrayToken[1]
	}

	return ""
}

func ValidateBearerToken(ctx context.Context, token string) (*ClaimToken, error) {
	cfg := configs.Get()
	if !strings.Contains(token, "Bearer") {
		return nil, errConstants.ErrUnautorized
	}

	tokenString := extractBearerToken(token)
	if tokenString == "" {
		return nil, errConstants.ErrUnautorized
	}

	claimToken := &ClaimToken{}
	tokenJwt, err := jwt.ParseWithClaims(tokenString, claimToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errConstants.ErrInvalidToken
		}

		jwtSecret := []byte(cfg.JwtConfig.JwtSecretKey)

		return jwtSecret, nil
	})

	if err != nil || !tokenJwt.Valid {
		return nil, errConstants.ErrUnautorized
	}

	return claimToken, nil
}
