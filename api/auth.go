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

// Login returns jwt token
// @Summary Login
// @Description Checks user credentials and returns JWT if ok
// @ID login
// @Tags user
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "User name"
// @Param password formData string true "Password"
// @Success 200 {object} tokenResponse
// @Failure 401 {object} errorResponse
// @Router /user/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := user{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	valid, err := auth.CheckСredentials(u.Username, u.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if !valid {
		writeError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	claims := &Claims{
		Username: u.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(720 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeOk(w, tokenResponse{
		Status: ok,
		Token:  tokenString,
	})
}

// SignUp creates new user
// @Summary Create new user
// @Description Creates new user
// @ID sign-up
// @Tags user
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "User name"
// @Param password formData string true "Password"
// @Success 200 {object} tokenResponse
// @Failure 401 {object} errorResponse
// @Router /user/signup [post]
func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := user{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	saved, err := auth.SaveСredentials(u.Username, u.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if !saved {
		writeError(w, http.StatusBadRequest, "Could not save password")
		return
	}

	claims := &Claims{
		Username: u.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(720 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeOk(w, tokenResponse{
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
			writeError(w, http.StatusUnauthorized, "No Authorization header found")
			return
		}
		authHeaderParts := strings.Fields(authHeader)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			writeError(w, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			return
		}
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(authHeaderParts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if !tkn.Valid {
			writeError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		r.Header.Add("user", claims.Username)
		next.ServeHTTP(w, r)
	})
}

// ValidateToken checks if token is valid
// @Summary Validate token
// @Description get string by ID
// @ID validate-token
// @Produce  json
// @Success 200 {object} tokenResponse
// @Failure 401 {object} errorResponse
// @Tags security
// @Security ApiKeyAuth
// @Router /token/validate [get]
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	writeOk(w, struct {
		Status responseStatus `json:"status"`
	}{Status: ok})
}
