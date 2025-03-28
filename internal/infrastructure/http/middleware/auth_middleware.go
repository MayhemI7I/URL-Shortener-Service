package middleware

import (
	"context"
	"net/http"
	


	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/handlers/urlhandler"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/pkg/utils/httputil"
	"github.com/MayhemI7I/URL-Shortener-Service/pkg/utils/jwtutil"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/auth"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/interfaces"
	"go.uber.org/zap"
)
// Middleware authenticates requests using JWT tokens and attaches user ID to context
// Auth is a middleware function that handles authentication.
// It extracts the access token from the request, validates it, and sets the user ID in the request context.
// If the access token is expired, it attempts to refresh it.
func Auth(s interfaces.AuthService) func(http.Handler) http.Handler {
   return func(next http.Handler) http.Handler {
   	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
   		// Extract the access token from the request.
   		accessToken, err := httputil.ExtractCookie(r, httputil.AccessTokenCookie)
   		if err != nil {
   			// Log the error and return an unauthorized error.
   			logger.Log.Debug("missing or invalid access token", zap.Error(err))
   			http.Error(w, "Unauthorized: missing or invalid access token", http.StatusUnauthorized)
   			return
   		}

   		// Parse the JWT token.
   		claims, err := jwtutil.ParseJWT(accessToken)
   		if err == nil {
   			// If the token is valid, set the user ID in the request context and call the next handler.
   			if claims.UserID == "" {
   				logger.Log.Debug("invalid claims: user ID is empty")
   				http.Error(w, "Unauthorized: invalid claims", http.StatusUnauthorized)
   				return
   			}
   			ctx := context.WithValue(r.Context(), urlhandler.UserIDKey, claims.UserID)
   			next.ServeHTTP(w, r.WithContext(ctx))
   			return
   		}

   		// If the token is expired, attempt to refresh it.
   		if auth.IsTokenExpired(err) {
   			 newAccessToken,err := auth.HandleTokenRefresh(w, r, s) 
			 if err != nil {
   				httputil.RespondWithError(w, err)
   				return
   			}

   			// Parse the new JWT token.
   			newClaims, err := jwtutil.ParseJWT(newAccessToken)
   			if err != nil {
   				logger.Log.Error("failed to parse refreshed token", zap.Error(err))
   				http.Error(w, "Internal server error", http.StatusInternalServerError)
   				return
   			}

   			// If the new token is valid, set the user ID in the request context and call the next handler.
   			if newClaims.UserID == "" {
   				logger.Log.Debug("invalid refreshed claims: user ID is empty")
   				http.Error(w, "Unauthorized: invalid claims", http.StatusUnauthorized)
   				return
   			}
   			ctx := context.WithValue(r.Context(), urlhandler.UserIDKey, newClaims.UserID)
   			next.ServeHTTP(w, r.WithContext(ctx))
   			return
   		}

   		// If the token is invalid, log the error and return an unauthorized error.
   		logger.Log.Debug("invalid access token", zap.Error(err))
   		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
   	})
   }
}


