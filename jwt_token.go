package main

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtTOken(key []byte) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)
	s, err := t.SignedString(key)
	if err != nil {
		return "", err
	}

	return s, nil
}

func VerifyJwtToken(tokenString string) (*jwt.Token, error) {
	tokenJwt, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		},
	)
	if err != nil {
		return &jwt.Token{}, err
	}

	return tokenJwt, nil
}
