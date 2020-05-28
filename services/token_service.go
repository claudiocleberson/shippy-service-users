package services

import (
	"time"

	"github.com/claudiocleberson/shippy-service-users/models"
	"github.com/dgrijalva/jwt-go"
)

var (
	key = []byte("mysupersecurekey")
)

type CustomClaims struct {
	User *models.User
	jwt.StandardClaims
}

type TokenService interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *models.User) (string, error)
}

func NewTokenService() TokenService {
	return &tokenService{}
}

type tokenService struct {
}

func (t *tokenService) Decode(tokenString string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (t *tokenService) Encode(user *models.User) (string, error) {

	expireToke := time.Now().Add(time.Hour * 72).Unix()

	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToke,
			Issuer:    "shippy-service-users",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}
