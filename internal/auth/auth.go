// auth/auth.go
package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/storage"
	"github.com/MayhemI7I/URL-Shortener-Service/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/utils/httputil"
)

func IsTokenExpired(err error) bool {
	if ve, ok := err.(*jwt.ValidationError); ok {
		return ve.Errors&jwt.ValidationErrorExpired != 0
	}
	return false
}

// HandleTokenRefresh handles the token refresh process.
// It extracts the refresh token from the request, gets a new access token from the storage,
// and sets the new access token and refresh token as cookies in the response.
func HandleTokenRefresh(w http.ResponseWriter, r *http.Request, s storage.Storage) (string, error) {
	// Extract the refresh token from the request.
	refreshToken, err := httputil.ExtractCookie(r, httputil.RefreshTokenCookie)
	if err != nil && err != http.ErrNoCookie {
		// Log the error and return an unauthorized error.
		logger.Log.Debug("missing or invalid refresh token", zap.Error(err))
		return "", fmt.Errorf("unauthorized: %v", err)
	}

	// Create a context with a timeout of 5 seconds.
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Get a new access token from the storage.
	newAccessToken, refreshToken, err := s.GetNewAccessToken(ctx, refreshToken)
	if err != nil {
		// Return the error if there was a problem getting the new access token.
		return "", err
	}

	// Set the new access token and refresh token as cookies in the response.
	httputil.SetCookie(w, httputil.AccessTokenCookie, newAccessToken, httputil.AccessTokenExpiry)
	httputil.SetCookie(w, httputil.RefreshTokenCookie, refreshToken, httputil.RefreshTokenExpiry)
	return newAccessToken, nil
}

func SetCSRFToken(w http.ResponseWriter, userID string) (string, error) {
	token := "csrf-token-example" // Заглушка
	w.Header().Set("X-CSRF-Token", token)
	return token, nil
}

func ValidateCSRFToken(r *http.Request) bool {
	return true // Заглушка,
}
