package authhandler

import (
	"net/http"
	"os"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/domain"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/storage"
	"github.com/MayhemI7I/URL-Shortener-Service/utils/jwtutil"
	"github.com/MayhemI7I/URL-Shortener-Service/utils/httputil"
	"github.com/MayhemI7I/URL-Shortener-Service/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthHandler struct{
	storage storage.Storage
}

func NewAuthHandler(s storage.Storage)*AuthHandler{
	return &AuthHandler{storage: s,}
}
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
   ctx := r.Context()
   id := uuid.New().String()
   refreshTokenStr, err := jwtutil.GenerateRefreshToken()
   if err != nil {
   	http.Error(w, err.Error(), http.StatusUnauthorized)
   	return
   }
   accessToken, err := jwtutil.GenerateAccessToken(id, os.Getenv("JWT_SECRET"))
   if err != nil {
   	http.Error(w, err.Error(), http.StatusUnauthorized)
   	return
   }
   
   // Создаем объект RefreshToken
   refreshToken := &domain.RefreshToken{
      User: domain.User{ID: id},
      Token: refreshTokenStr,
      ExpiresAt: time.Now().Add(jwtutil.RefreshTokenExpiration),
   }
   
   err = ah.storage.SaveRefreshToken(ctx, refreshToken)
   if err != nil {
   	http.Error(w, err.Error(), http.StatusUnauthorized)
   	return
   }
   httputil.SetCookie(w, httputil.AccessTokenCookie, accessToken, httputil.AccessTokenExpiry)
   httputil.SetCookie(w, httputil.RefreshTokenCookie, refreshTokenStr, httputil.RefreshTokenExpiry)
   w.WriteHeader(http.StatusOK)
   logger.Log.Info("successfully registered user", zap.String("userID", id))
}

