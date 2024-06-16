package jwtutil

import (
	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JWTUtilRepo struct {
	App config.Config
}

var Repo = JWTUtilRepo{}

func NewJWTUtil(a config.Config) {
	Repo.App = a
}

func (j *JWTUtilRepo) Generate(claims jwt.MapClaims) (string, error) {
	// Crate JWT to authenticate user
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //jwt.SigningMethodES256
	jwt, err := t.SignedString(j.App.GetJWTSecretKey())
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (j *JWTUtilRepo) Parse(token string) (jwt.MapClaims, error) {

	errInvalidToken := status.Errorf(codes.Unauthenticated, "invalid token")

	// Parse Token
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return j.App.GetJWTSecretKey(), nil
	})
	if err != nil {
		return nil, errInvalidToken
	}

	// Get claims
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errInvalidToken
	}

	return claims, nil
}
