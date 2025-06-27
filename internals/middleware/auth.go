package middlewares

import (
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
)

func AuthMiddleware(ah http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err := utils.ValidateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ah.ServeHTTP(w, r)
	})
}
