package utility

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserInfo struct {
	ExpiresAt int64
	UserId    string
	RoleID    string
	SecretKey []byte
}

func (userClaim *UserInfo) CreateNewToken() (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "movy theatre",
		"sub":  userClaim.UserId,
		"role": userClaim.RoleID,
		"exp":  userClaim.ExpiresAt,
		"iat":  time.Now().Unix(),
	})

	token, err := claims.SignedString(userClaim.SecretKey)

	if err != nil {
		return "", err
	}
	return token, nil
}

func (userClaim *UserInfo) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return userClaim.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid error token provided")
	}

	if iss, ok := claims["iss"]; !ok || iss != "movy theatre" {
		return nil, errors.New("invalid issuer")
	}

	return claims, err
}
