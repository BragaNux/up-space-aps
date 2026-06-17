package middleware

import (
	"net/http"
	"strings"
)

// WithCORS libera o acesso pra lista de origens permitidas (vem da config, separadas por virgula).
// Se a origem da requisicao bater com uma da lista, devolve ela mesma; senao usa a primeira como fallback (util pra curl/clientes sem browser)
func WithCORS(allowedOrigins string) func(http.Handler) http.Handler {
	origins := strings.Split(allowedOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestOrigin := r.Header.Get("Origin")
			allowOrigin := origins[0]

			for _, o := range origins {
				if o == "*" || o == requestOrigin {
					allowOrigin = requestOrigin
					break
				}
			}

			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
