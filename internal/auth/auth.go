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



// Helper functions

func IsTokenExpired(err error) bool {
	if ve, ok := err.(*jwt.ValidationError); ok {
		return ve.Errors&jwt.ValidationErrorExpired != 0
	}
	return false
}

func HandleTokenRefresh(w http.ResponseWriter, r *http.Request, s storage.Storage) error {
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

func SetCSRFToken(w http.ResponseWriter, userID string) (string, error) {
	token := "csrf-token-example" // Заглушка
	w.Header().Set("X-CSRF-Token", token)
	return token, nil
}

func ValidateCSRFToken(r *http.Request) bool {
	return true // Заглушка, 
}
