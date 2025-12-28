package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/inspection-tool/backend/internal/domain"
)

// Package middleware provides HTTP middleware for authentication,
// authorization and logging. Middlewares attach authenticated user
// information to the request context and provide helpers to write API
// responses in a consistent JSON format.

// AuthContextKey is used to store the authenticated user in context
type AuthContextKey string

const UserContextKey AuthContextKey = "user"

// AuthMiddleware validates JWT tokens and adds user info to context.
// The middleware expects an `Authorization: Bearer <token>` header and
// uses `jwtSecret` to validate token signatures. On success it extracts
// commonly used claims into a `domain.User` struct and stores it under
// `UserContextKey` for downstream handlers to consume.
func AuthMiddleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header", nil)
				return
			}

			// Parse "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization header format", nil)
				return
			}

			tokenString := parts[1]

			// Parse and validate JWT
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token", nil)
				return
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token claims", nil)
				return
			}

			// Create user from claims
			user := &domain.User{
				UserID:   claims["sub"].(string),
				Username: claims["cognito:username"].(string),
				Email:    claims["email"].(string),
				Role:     domain.Role(claims["custom:role"].(string)),
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GuestAuthMiddleware validates guest session tokens. Guest tokens are
// expected to contain minimal claims and will produce a `domain.User` with
// `Role` set to `domain.RoleGuestViewer`. This middleware is suitable for
// endpoints that allow temporary guest access (e.g., viewer links).
func GuestAuthMiddleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header", nil)
				return
			}

			// Parse "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization header format", nil)
				return
			}

			tokenString := parts[1]

			// Parse and validate JWT
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token", nil)
				return
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token claims", nil)
				return
			}

			// Create guest user from claims
			user := &domain.User{
				UserID: claims["sub"].(string),
				Role:   domain.RoleGuestViewer,
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole middleware ensures the authenticated user has at least one
// of the provided roles. If the user is missing or does not have the
// required role, it returns a 403 Forbidden response.
func RequireRole(allowedRoles ...domain.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(UserContextKey).(*domain.User)

			for _, role := range allowedRoles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			writeErrorResponse(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions", nil)
		})
	}
}

// LoggingMiddleware logs incoming HTTP requests using structured
// `slog.Logger`. It records method, path and remote address. Keep this
// middleware early in the chain so that all requests are logged.
func LoggingMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("HTTP request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
			)
			next.ServeHTTP(w, r)
		})
	}
}

// ErrorResponse represents a standard error response sent to clients.
// The `Code` field is a machine-friendly error code (see domain errors)
// while `Message` provides a human-friendly explanation. `Details` can
// include structured context useful for debugging.
type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// SuccessResponse represents a standard success response envelope.
// Use `WriteSuccessResponse` in handlers to return data in a
// predictable structure that clients can parse.
type SuccessResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// writeErrorResponse writes an error response with the given HTTP
// status code and JSON body. It is kept unexported to centralize
// formatting logic while exported wrappers are provided below.
func writeErrorResponse(w http.ResponseWriter, statusCode int, code, message string, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}

// WriteErrorResponse is exported for use in handlers when an error
// needs to be returned to the client.
func WriteErrorResponse(w http.ResponseWriter, statusCode int, code, message string, details map[string]interface{}) {
	writeErrorResponse(w, statusCode, code, message, details)
}

// WriteSuccessResponse writes a standardized success envelope containing
// the provided `data` payload.
func WriteSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{
		Code:    "SUCCESS",
		Message: "Operation successful",
		Data:    data,
	})
}
