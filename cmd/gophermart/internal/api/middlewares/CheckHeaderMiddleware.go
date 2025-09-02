package middlewares

import (
	"fmt"
	"net/http"
)

func CheckHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			fmt.Println(r.Header)
			return
		}
		next.ServeHTTP(w, r)
	})
}
