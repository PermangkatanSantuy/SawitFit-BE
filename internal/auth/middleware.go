package auth

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Middleware(verifier *Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("Authorization")

			if header == "" {
				http.Error(w, "missing auth header.", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(header, "Bearer ")

			token, err := verifier.Verify(tokenString)
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)

			user := &UserClaims {
				Sub: claims["sub"].(string),
				Email: claims["email"].(string),
				Role: claims["role"].(string),
			}

			ctx := WithUser(r.Context(), user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}