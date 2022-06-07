package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyClaims struct {
	Username string
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2

var MySecretKey = []byte("cisco")

func GenToken(username string) (string, error) {
	c := MyClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "kubernetes",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString(MySecretKey)

}

func ParseToken(tokenString string) (*MyClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
