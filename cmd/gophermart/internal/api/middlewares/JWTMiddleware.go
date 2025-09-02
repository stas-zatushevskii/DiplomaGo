package middlewares

import (
	"context"
	"fmt"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/service"
	"net/http"
)

func JWTMiddleware(service *service.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("JWT")
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				fmt.Println("---------------------------------------------")
				fmt.Println("COOKIE NOT FOUND")
				fmt.Println("---------------------------------------------")
				return
			}
			userID, err := service.UserService.GetUserIDFromJwt(cookie.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				fmt.Println("---------------------------------------------")
				fmt.Println("BAD COOKIE")
				fmt.Println("---------------------------------------------")
				return
			}

			fmt.Println("---------------------------------------------")
			fmt.Println("SET USER ID IN CONTEXT:", userID)
			fmt.Println("---------------------------------------------")
			ctx := context.WithValue(r.Context(), "UserID", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
