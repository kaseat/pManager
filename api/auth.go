package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kaseat/pManager/auth"
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

	u := user{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	_, err := auth.CheckСredentials(u.Username, u.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	claims := &Claims{
		Username: u.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(50 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeOk(w, struct {
		Status responseStatus `json:"status"`
		Token  string         `json:"token"`
	}{
		Status: ok,
		Token:  tokenString,
	})
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
