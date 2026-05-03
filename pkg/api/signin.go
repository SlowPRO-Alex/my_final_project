package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignInRequest struct {
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

var jwtSecret = []byte("some-secret-key")

func someHash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func generateJWT(pass string) (string, error) {
	claims := jwt.MapClaims{
		"hash": someHash(pass),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "incorrect request", http.StatusBadRequest)
		return
	}
	authPass := req.Password
	fmt.Printf("input password: %s\n", authPass)
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		writeJson(w, SignInResponse{Error: "auth off"})
		return
	}
	if authPass != pass {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}
	token, err := generateJWT(pass)
	if err != nil {
		writeJson(w, SignInResponse{Error: "error generate token"})
		return
	}
	fmt.Printf("correct password: %s, JWT-token: %s\n", pass, token)
	writeJson(w, SignInResponse{Token: token})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var tokenString string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			} else {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// проверяем метод подписи
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "authentification required", http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			hashInToken, ok := claims["hash"].(string)
			if !ok {
				http.Error(w, "invalid token data", http.StatusUnauthorized)
				return
			}

			if hashInToken != someHash(pass) {
				http.Error(w, "invalid token hash", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}
