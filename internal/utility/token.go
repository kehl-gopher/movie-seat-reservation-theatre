package utility

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaim struct {
	ExpiresAt int64
	UserId    string
	Role      uint8
	SecretKey []byte
}

func (userClaim *AccessTokenClaim) CreateNewToken() (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "movy theatre",
		"sub":  userClaim.UserId,
		"role": userClaim.Role,
		"exp":  time.Now().Add(time.Duration(userClaim.ExpiresAt) * time.Minute).Unix(),
		"iat":  time.Now().Unix(),
	})

	token, err := claims.SignedString(userClaim.SecretKey)

	if err != nil {
		return "", err
	}
	return token, nil
}

func (userClaim *AccessTokenClaim) ValidateToken(tokenString string) (jwt.MapClaims, error) {
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
