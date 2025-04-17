package utility

import (
	"errors"
	"fmt"
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
		"iss":     "movy theatre",
		"user_id": userClaim.UserId,
		"role":    userClaim.Role,
		"exp":     time.Now().AddDate(0, 0, int(userClaim.ExpiresAt)).Unix(),
	})

	token, err := claims.SignedString(userClaim.SecretKey)

	if err != nil {
		return "", err
	}
	return token, nil
}

func (userClaim *AccessTokenClaim) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method %v", t.Header["alg"]))
		}
		return userClaim.SecretKey, nil
	})

	if err != nil {
		return nil, errors.New("invalid token provided")
	}

	return token, err
}

func (userClaim *AccessTokenClaim) ExtractClaims(tokenString string) (string, uint8, error) {
	token, err := userClaim.ValidateToken(tokenString)
	if err != nil {
		return "", 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["user_id"].(string)
		role := uint8(claims["role"].(float64))
		fmt.Println("---------->")
		return userId, uint8(role), nil
	}
	return "", 0, errors.New("invalid token provided")
}
