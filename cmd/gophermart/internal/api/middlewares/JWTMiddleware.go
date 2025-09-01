package middlewares

import (
	"context"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/service"
	"net/http"
)

func JWTMiddleware(service *service.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("JWT")
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			userID, err := service.UserService.GetUserIdFromJwt(cookie.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
