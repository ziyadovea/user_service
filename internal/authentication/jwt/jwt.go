package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type accessTokenClaims struct {
	UserID int64 `json:"sub"`
	jwt.StandardClaims
}

type refreshTokenClaims struct {
	UserID int64 `json:"sub"`
	jwt.StandardClaims
}

type authenticator struct {
	accessTokenSecretKey           []byte
	refreshTokenSecretKey          []byte
	accessTokenExpirationDuration  time.Duration
	refreshTokenExpirationDuration time.Duration
}

func NewAuthenticator(
	accessTokenSecretKey []byte,
	refreshTokenSecretKey []byte,
	accessTokenExpirationDuration time.Duration,
	refreshTokenExpirationDuration time.Duration,
) authenticator {
	return authenticator{
		accessTokenSecretKey:           accessTokenSecretKey,
		refreshTokenSecretKey:          refreshTokenSecretKey,
		accessTokenExpirationDuration:  accessTokenExpirationDuration,
		refreshTokenExpirationDuration: refreshTokenExpirationDuration,
	}
}

func (a authenticator) CreateAccessToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.accessTokenExpirationDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	tokenString, err := token.SignedString(a.accessTokenSecretKey)
	if err != nil {
		return "", fmt.Errorf("unable to signed token: %w")
	}

	return tokenString, nil
}

func (a authenticator) CreateRefreshToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.refreshTokenExpirationDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	tokenString, err := token.SignedString(a.refreshTokenSecretKey)
	if err != nil {
		return "", fmt.Errorf("unable to signed token: %w")
	}

	return tokenString, nil
}

func (a authenticator) VerifyAccessToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.accessTokenSecretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*accessTokenClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	return claims.UserID, nil
}

func (a authenticator) VerifyRefreshToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &refreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.refreshTokenSecretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*refreshTokenClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	return claims.UserID, nil
}
