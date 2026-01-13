package middleware

import (
	"context"
	"inventory-api/utils"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

type AuthContext struct {
	UserID   uint
	Username string
	Role     string
}

func HumaAuthMiddleware(api huma.API, jwtSecret string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authHeader := ctx.Header("Authorization")
		if authHeader == "" {
			huma.WriteErr(api, ctx, 401, "Authorization header required")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			huma.WriteErr(api, ctx, 401, "Invalid authorization format. Use: Bearer <token>")
			return
		}

		token := parts[1]
		claims, err := utils.ValidateToken(token, jwtSecret)
		if err != nil {
			huma.WriteErr(api, ctx, 401, "Invalid or expired token")
			return
		}

		// Store auth info in context
		authCtx := &AuthContext{
			UserID:   claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		// Create new context with auth info
		newCtx := context.WithValue(ctx.Context(), "auth", authCtx)
		ctx = huma.WithContext(ctx, newCtx)

		next(ctx)
	}
}

func GetAuthContext(ctx context.Context) *AuthContext {
	if auth, ok := ctx.Value("auth").(*AuthContext); ok {
		return auth
	}
	return nil
}
