package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	accessSecretKey  string
	refreshSecretKey string
}

func NewJWTMaker(accessSecretKey string, refreshSecretKey string) *JWTMaker {
	return &JWTMaker{
		accessSecretKey:  accessSecretKey,
		refreshSecretKey: refreshSecretKey,
	}
}

func (maker *JWTMaker) CreateRefreshToken(id int, email string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, email, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(maker.refreshSecretKey))
	if err != nil {
		return "", nil, fmt.Errorf("error sighning token: %w", err)
	}

	return tokenStr, claims, nil
}

func (maker *JWTMaker) CreateAccessToken(id int, email string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, email, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(maker.accessSecretKey))
	if err != nil {
		return "", nil, fmt.Errorf("error sighning token: %w", err)
	}

	return tokenStr, claims, nil
}

func (maker *JWTMaker) VerifyAccessToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(maker.accessSecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if ok != ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil

}

func (maker *JWTMaker) VerifyRefreshToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(maker.refreshSecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if ok != ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil

}
