package middlewares

import (
	"log"
	"net/http"

	"github.com/VsenseTechnologies/biometric_http_server/pkg/utils"
)

func AuthMiddleware(ah http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Println("token is empty")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err := utils.ValidateToken(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ah.ServeHTTP(w, r)
	})
}
