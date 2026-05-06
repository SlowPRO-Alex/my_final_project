package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthPass struct {
	Password string `json:"password"`
}

var jwtSecret = []byte("some-secret-key")
var pass = os.Getenv("TODO_PASSWORD")

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
	var authPass AuthPass
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &authPass)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	//authPass := r.URL.Query().Get("password")
	log.Printf("TODO_PASSWORD: %s\n", pass)
	log.Printf("input password: %s\n", authPass.Password)
	if pass == "" {
		writeJson(w, map[string]string{"error": "auth off"}, http.StatusBadRequest)
		return
	}
	if authPass.Password != pass {
		writeJson(w, map[string]string{"error": "invalid password"}, http.StatusBadRequest)
		return
	}
	token, err := generateJWT(pass)
	if err != nil {
		writeJson(w, map[string]string{"error": "error generate token"}, http.StatusInternalServerError)
		return
	}
	log.Printf("correct password: %s, JWT-token: %s\n", pass, token)
	writeJson(w, map[string]string{"token": token}, http.StatusOK)
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		if len(pass) > 0 {
			var tokenString string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			} else {
				writeJson(w, map[string]string{"error": "Authentication required"}, http.StatusUnauthorized)
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
				writeJson(w, map[string]string{"error": "Authentication required"}, http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeJson(w, map[string]string{"error": "invalid token claims"}, http.StatusUnauthorized)
				return
			}

			hashInToken, ok := claims["hash"].(string)
			if !ok {
				writeJson(w, map[string]string{"error": "invalid token data"}, http.StatusUnauthorized)
				return
			}

			if hashInToken != someHash(pass) {
				writeJson(w, map[string]string{"error": "invalid token hash"}, http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}
