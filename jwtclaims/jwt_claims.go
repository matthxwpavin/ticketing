package jwtclaims

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/matthxwpavin/ticketing/env"
)

type jwtContextKey struct{}

var jwtKey jwtContextKey

func WithContext(parent context.Context, claims *CustomClaims) context.Context {
	return context.WithValue(parent, jwtKey, claims)
}

func FromContext(ctx context.Context) *CustomClaims {
	claims, _ := ctx.Value(jwtKey).(*CustomClaims)
	return claims
}

type CustomClaims struct {
	Metadata Metadata `json:"metadata"`
	jwt.RegisteredClaims
}

type Metadata struct {
	Email  string `json:"email"`
	UserID string `json:"userID"`
}

func IssueToken(data Metadata) (string, error) {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			Metadata: data,
			RegisteredClaims: jwt.RegisteredClaims{
				// A usual scenario is to set the expiration time relative to the current time
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				// Issuer:    "test",
				// Subject:   "somebody",
				// ID: "some_id",
				// Audience:  []string{"somebody_else"},
			},
		},
	).SignedString([]byte(jwtKeyFromEnv()))
}

func Parse(jwtStr string) (*CustomClaims, error) {
	_, err := jwt.Parse(jwtStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKeyFromEnv()), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	rawClaims := strings.Split(jwtStr, ".")[1]
	decoded, err := base64.RawStdEncoding.DecodeString(rawClaims)
	if err != nil {
		return nil, err
	}

	claims := new(CustomClaims)
	return claims, json.Unmarshal(decoded, claims)
}

func jwtKeyFromEnv() string {
	return env.JwtSecret.Value()
}
