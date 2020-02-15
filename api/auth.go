package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte("my_secret_key")

// Claims represents users claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GetToken returns jwt token
func GetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := &Claims{
		Username: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := struct {
		Status responseStatus
		Token  string
	}{
		Status: ok,
		Token:  tokenString,
	}

	bytes, err := json.Marshal(&resp)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// VerifyTokenMiddleware verifies token, if token ok, allawes request pass-through
func VerifyTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, http.StatusBadRequest, "No Authorization header found")
			return
		}
		authHeaderParts := strings.Fields(authHeader)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			writeError(w, http.StatusBadRequest, "Authorization header format must be Bearer {token}")
			return
		}
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(authHeaderParts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if !tkn.Valid {
			writeError(w, http.StatusBadRequest, "Invalid token")
			return
		}

		r.Header.Add("user", claims.Username)
		next.ServeHTTP(w, r)
	})
}
