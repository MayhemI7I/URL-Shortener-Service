package middleware

import (
	"context"
	"net/http"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/storage"
	"github.com/MayhemI7I/URL-Shortener-Service/handlers/urlhandler"
	"github.com/MayhemI7I/URL-Shortener-Service/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/utils/httputil"
	"github.com/MayhemI7I/URL-Shortener-Service/utils/jwtutil"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/auth"
	"go.uber.org/zap"
)
// Middleware authenticates requests using JWT tokens and attaches user ID to context
func Auth(s storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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
				ctx := context.WithValue(r.Context(), urlhandler.UserIDKey, claims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if auth.IsTokenExpired(err) {
				if err := auth.HandleTokenRefresh(w, r, s); err != nil {
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
				ctx := context.WithValue(r.Context(), urlhandler.UserIDKey, newClaims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			logger.Log.Debug("invalid access token", zap.Error(err))
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		})
	}
}
