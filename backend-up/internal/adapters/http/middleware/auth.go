package middleware

import (
	"context"
	"net/http"
	"strings"

	"up-espaco/backend/internal/adapters/security"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
	roleKey   contextKey = "role"
)

// extractToken pega o token do header "Authorization: Bearer xxx"
func extractToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}

// RequireAuth bloqueia a rota com 401 se nao vier um token valido
func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				http.Error(w, `{"Error":"token de autorização ausente"}`, http.StatusUnauthorized)
				return
			}

			userID, role, err := security.ParseToken(secret, token)
			if err != nil {
				http.Error(w, `{"Error":"token inválido ou expirado"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, roleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth coloca o usuario no contexto se o token vier valido, mas nunca bloqueia a rota
func OptionalAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token != "" {
				if userID, role, err := security.ParseToken(secret, token); err == nil {
					ctx := context.WithValue(r.Context(), userIDKey, userID)
					ctx = context.WithValue(ctx, roleKey, role)
					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// UserIDFromContext pega o id do usuario logado que o middleware de auth guardou no contexto
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}

// RoleFromContext pega o papel (profissional/responsavel) do usuario logado
func RoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey).(string)
	return role, ok
}
