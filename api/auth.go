package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte("my_secret_key")

// GetToken returns jwt token
func GetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user"] = "admin"
	claims["exp"] = time.Now().Add(time.Minute).Unix()
	tokenString, _ := token.SignedString(secret)
	fmt.Println(tokenString)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
}

// VerifyToken verfies jwt token
func VerifyToken(r *http.Request) (bool, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, errors.New("No Authorization header found")
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		return false, errors.New("Authorization header format must be Bearer {token}")
	}
	claims := &jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(authHeaderParts[1], claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return false, err
	}
	fmt.Println(*claims)

	return tkn.Valid, nil
}
