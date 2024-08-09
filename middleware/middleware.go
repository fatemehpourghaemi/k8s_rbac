package middleware

import (
	"k8s_rbac/utilities"
	"fmt"
	"net/http"
	"strings"
)

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token != strings.TrimSpace(utilities.Token) {
			fmt.Printf(r.Header.Get("Authorization"))

			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	}
}
