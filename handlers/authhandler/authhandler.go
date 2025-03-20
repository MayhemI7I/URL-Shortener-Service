package authhandler

import (
	"net/http"
	"os"
	"time"

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
   refreshToken, err := jwtutil.GenerateRefreshToken()
   if err != nil {
   	http.Error(w, err.Error(), http.StatusUnauthorized)
   	return
   }
   accessToken, err := jwtutil.GenerateAccessToken(id, os.Getenv("JWT_SECRET"))
   if err != nil {
   	http.Error(w, err.Error(), http.StatusUnauthorized)
   	return
   }
   err = ah.storage.SaveRefreshToken(ctx, refreshToken, id, time.Now().Add(jwtutil.RefreshTokenExpiration))
   if err != nil {
   	http.Error(w, err.Error(), http.StatusUnauthorized)
   	return
   }
   httputil.SetCookie(w, httputil.AccessTokenCookie, accessToken, httputil.AccessTokenExpiry)
   httputil.SetCookie(w, httputil.RefreshTokenCookie, refreshToken, httputil.RefreshTokenExpiry)
   w.WriteHeader(http.StatusOK)
   logger.Log.Info("successfully registered user", zap.String("userID", id))
}

