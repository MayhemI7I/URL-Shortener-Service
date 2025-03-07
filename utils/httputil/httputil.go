package httputil

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"context"

	"go.uber.org/zap"

	"github.com/MayhemI7I/URL-Shortener-Service/domain"
	"github.com/MayhemI7I/URL-Shortener-Service/utils/jwtutil"
	"github.com/MayhemI7I/URL-Shortener-Service/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/utils/token"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refreshToken" 
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 24 * time.Hour
)

// SetJWTCookieWithNewRefresh sets access and refresh token cookies with a new refresh token
func SetJWTCookieWithNewRefresh(w http.ResponseWriter, userID, secretKey string) error {
	accessToken, err := jwtutil.GenerateAccessToken(userID, secretKey)
	if err != nil {
		return err
	}
	refreshToken, err := token.GenerateRefreshToken(userID)
	if err != nil {
		return err
	}

	SetCookie(w, AccessTokenCookie, accessToken, AccessTokenExpiry)
	SetCookie(w, RefreshTokenCookie, refreshToken, RefreshTokenExpiry)
	return nil
}

// SetJWTCookieWithOldRefresh sets access and refresh token cookies with an existing refresh token
func SetJWTCookieWithOldRefresh(w http.ResponseWriter, userID, secretKey string) error {
	accessToken, err := jwtutil.GenerateAccessToken(userID, secretKey)
	if err != nil {
		return err
	}
	
	SetCookie(w, AccessTokenCookie, accessToken, AccessTokenExpiry)
	
	return nil
}

// ExtractCookie retrieves a cookie value by name from the request
func ExtractCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil || cookie.Value == "" {
		return "", fmt.Errorf("empty or missing %s cookie", name)
	}
	return cookie.Value, nil
}

// SetCookie sets a secure cookie with the given name, value, and expiry
func SetCookie(w http.ResponseWriter, name, value string, expiry time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode, // Можно сделать Strict в зависимости от требований
		Expires:  time.Now().Add(expiry),
	})
}

// RespondWithError sends an appropriate HTTP error response based on the error type
func RespondWithError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrURLNotFound):
		http.Error(w, "URL not found", http.StatusNotFound)
	case errors.Is(err, domain.ErrURLExists):
		http.Error(w, "URL already exists", http.StatusConflict)
	case errors.Is(err, domain.ErrTokenNotFound):
		http.Error(w, "Unauthorized: refresh token not found", http.StatusUnauthorized)
	case errors.Is(err, domain.ErrTokenExpired):
		http.Error(w, "Unauthorized: refresh token expired", http.StatusUnauthorized)
	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
	default:
		logger.Log.Error("internal error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
