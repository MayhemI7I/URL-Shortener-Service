// auth/auth.go
package auth

import (
	"context"
	"net/http"
	"time"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"


	"local/utils/httputil"
	"local/internal/storage"
	"local/utils/jwtutil"
	"local/logger"
)

// Middleware authenticates requests using JWT tokens and attaches user ID to context
func Middleware(s storage.Storage, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := httputil.ExtractCookie(r, httputil.AccessTokenCookie)
		if err != nil {
			logger.Log.Debug("missing or invalid access token", zap.Error(err))
			http.Error(w, "Unauthorized: missing or invalid access token", http.StatusUnauthorized)
			return
		}

		claims, err := jwtutil.ParseJWT(accessToken)
		if err == nil {
			if claims.UserID == "" {
				logger.Log.Debug("invalid claims: user ID is empty")
				http.Error(w, "Unauthorized: invalid claims", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			if _, err := setCSRFToken(w, claims.UserID); err != nil {
				logger.Log.Error("failed to set CSRF token", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if isTokenExpired(err) {
			if err := handleTokenRefresh(w, r, s); err != nil {
				httputil.RespondWithError(w, err)
				return
			}
			newAccessToken, err := httputil.ExtractCookie(r, httputil.AccessTokenCookie)
			if err != nil {
				logger.Log.Error("failed to extract refreshed access token", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			newClaims, err := jwtutil.ParseJWT(newAccessToken)
			if err != nil {
				logger.Log.Error("failed to parse refreshed token", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if newClaims.UserID == "" {
				logger.Log.Debug("invalid refreshed claims: user ID is empty")
				http.Error(w, "Unauthorized: invalid claims", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "userID", newClaims.UserID)
			if _, err := setCSRFToken(w, newClaims.UserID); err != nil {
				logger.Log.Error("failed to set CSRF token", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		logger.Log.Debug("invalid access token", zap.Error(err))
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
	})
}

// Helper functions

func isTokenExpired(err error) bool {
	if ve, ok := err.(*jwt.ValidationError); ok {
		return ve.Errors&jwt.ValidationErrorExpired != 0
	}
	return false
}

func handleTokenRefresh(w http.ResponseWriter, r *http.Request, s storage.Storage) error {
	refreshToken, err := httputil.ExtractCookie(r, httputil.RefreshTokenCookie)
	if err != nil {
		logger.Log.Debug("missing or invalid refresh token", zap.Error(err))
		return fmt.Errorf("unauthorized: %v", err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	newAccessToken, newRefreshToken, err := s.GetNewAccessToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	httputil.SetCookie(w, httputil.AccessTokenCookie, newAccessToken, httputil.AccessTokenExpiry)
	httputil.SetCookie(w, httputil.RefreshTokenCookie, newRefreshToken, httputil.RefreshTokenExpiry)
	return nil
}

func setCSRFToken(w http.ResponseWriter, userID string) (string, error) {
	token := "csrf-token-example" // Заглушка
	w.Header().Set("X-CSRF-Token", token)
	return token, nil
}

func ValidateCSRFToken(r *http.Request) bool {
	return true // Заглушка, 
}