package utils

import (
	"errors"
	"time"

	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret     []byte
	expireTime time.Duration
}

func NewJWTService(secret string, expireHours int) *JWTService {
	return &JWTService{
		secret:     []byte(secret),
		expireTime: time.Duration(expireHours) * time.Hour,
	}
}

func (j *JWTService) GenerateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.Name,
		"role":    user.Role,
		"exp":     time.Now().Add(j.expireTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTService) ValidateToken(tokenString string) (*dto.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	name, ok := claims["name"].(string)
	if !ok {
		return nil, errors.New("invalid name in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role in token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid exp in token")
	}

	return &dto.JWTClaims{
		UserID: int(userID),
		Name:   name,
		Role:   role,
		Exp:    int64(exp),
	}, nil
}
