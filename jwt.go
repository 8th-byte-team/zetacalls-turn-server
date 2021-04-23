package main

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func ParseJWT(JWTSignature string, JWT string) (*jwt.StandardClaims, error) {
	res := jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(JWT, &res, func(token *jwt.Token) (interface{}, error) {
		return JWTSignature, nil
	})

	if err != nil {
		if err, ok := err.(jwt.ValidationError); ok {
			if err.Errors == jwt.ValidationErrorExpired {
				return &res, errors.New("token is expired")
			}
		}

		if strings.Contains(err.Error(), "is expired") {
			return &res, errors.New("token is expired")
		}

		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("bad token")
	}

	return &res, nil
}
