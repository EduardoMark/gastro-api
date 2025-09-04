package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/EduardoMark/gastro-api/internal/auth"
)

type JWTMiddleware struct {
	authService *auth.AuthJWTService
}

func NewJWTMiddleware(authService *auth.AuthJWTService) *JWTMiddleware {
	return &JWTMiddleware{
		authService: authService,
	}
}

type contentKey string

const CtxUserId contentKey = "user_id"
const CtxUserRole contentKey = "role"

func (m *JWTMiddleware) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "missing auth header",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid auth header",
			})
			return
		}

		tokenStr := parts[1]
		claims, err := m.authService.VerifyToken(tokenStr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserId, claims.UserID)
		ctx = context.WithValue(ctx, CtxUserRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
