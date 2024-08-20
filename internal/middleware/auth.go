package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/JoshElias/go-web-server/internal"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			fmt.Println("1")
			internal.RespondWithError(w, 401)
			return
		}
		userIdString, err := token.Claims.GetSubject()
		if err != nil {
			fmt.Println("2")
			internal.RespondWithError(w, 401)
			return
		}
		userId, err := strconv.Atoi(userIdString)
		if err != nil {
			fmt.Println("3")
			internal.RespondWithError(w, 401)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
