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

type contextKey string

const authContextKey contextKey = "auth"

// HumaAuthMiddleware validates JWT token and adds auth context
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
			huma.WriteErr(api, ctx, 401, "Invalid or expired token: "+err.Error())
			return
		}

		// Store auth info in context
		authCtx := &AuthContext{
			UserID:   claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		// Create new context with auth info
		newCtx := context.WithValue(ctx.Context(), authContextKey, authCtx)
		ctx = huma.WithContext(ctx, newCtx)

		next(ctx)
	}
}

// RequireRole creates a middleware that checks if user has required role
func RequireRole(api huma.API, jwtSecret string, allowedRoles ...string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// First authenticate
		HumaAuthMiddleware(api, jwtSecret)(ctx, func(authenticatedCtx huma.Context) {
			// Then check role
			auth := GetAuthContext(authenticatedCtx.Context())
			if auth == nil {
				huma.WriteErr(api, authenticatedCtx, 401, "Authentication required")
				return
			}

			// Check if user has one of the allowed roles
			hasRole := false
			for _, role := range allowedRoles {
				if auth.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				huma.WriteErr(api, authenticatedCtx, 403, "Insufficient permissions. Required role: "+strings.Join(allowedRoles, " or "))
				return
			}

			next(authenticatedCtx)
		})
	}
}

// RequireOwnerOrAdmin checks if user is the owner of resource or admin
func RequireOwnerOrAdmin(api huma.API, jwtSecret string, resourceUserID uint) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		HumaAuthMiddleware(api, jwtSecret)(ctx, func(authenticatedCtx huma.Context) {
			auth := GetAuthContext(authenticatedCtx.Context())
			if auth == nil {
				huma.WriteErr(api, authenticatedCtx, 401, "Authentication required")
				return
			}

			// Allow if user is admin or owns the resource
			if auth.Role != "admin" && auth.UserID != resourceUserID {
				huma.WriteErr(api, authenticatedCtx, 403, "Access denied. You can only access your own resources")
				return
			}

			next(authenticatedCtx)
		})
	}
}

// GetAuthContext retrieves auth context from context
func GetAuthContext(ctx context.Context) *AuthContext {
	if auth, ok := ctx.Value(authContextKey).(*AuthContext); ok {
		return auth
	}
	return nil
}

// IsAdmin checks if current user is admin
func IsAdmin(ctx context.Context) bool {
	auth := GetAuthContext(ctx)
	return auth != nil && auth.Role == "admin"
}

// IsOwner checks if current user owns the resource
func IsOwner(ctx context.Context, resourceUserID uint) bool {
	auth := GetAuthContext(ctx)
	return auth != nil && auth.UserID == resourceUserID
}

// IsOwnerOrAdmin checks if current user is owner or admin
func IsOwnerOrAdmin(ctx context.Context, resourceUserID uint) bool {
	return IsAdmin(ctx) || IsOwner(ctx, resourceUserID)
}
