package middlewares

import "net/http"

// func CorsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		// If it's a preflight request, just return 200 OK
// 		if r.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

var allowedOrigins = map[string]bool{
	"http://localhost:3000":                 true,
	"http://localhost:3001":                 true,
	"https://biometric.admin.vithsutra.com": true,
}

func CorsMiddleware(ah http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")

		if allowedOrigins[origin] {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", origin)                                           // Allow all origins
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")               // Allow specific methods
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With") // Allow headers
			w.Header().Set("Access-Control-Allow-Credentials", "true")                                      // Allow cookies (if needed)
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		ah.ServeHTTP(w, r)
	})
}
