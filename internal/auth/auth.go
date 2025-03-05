package auth

import (
	"crypto/rand"
	"encoding/hex"
	"local/config"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	
)



func SetJWTCookie(w http.ResponseWriter, token string, secretKey *config.Config)(error) {
	accsesToken, err := GenerateAccessToken(token, secretKey)
	if err != nil {
		return err
	}
	refreshToken, err := GenerateRefresToken(token)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name: "accses_token",
		Value: accsesToken,
		Path: "/",
		MaxAge: int(15 * time.Minute),
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w,  &http.Cookie{
		Name: "refreshToken",
		Value: refreshToken,
		Path: "/",
		MaxAge: int(24 * time.Hour),
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil

}


type Claims struct{
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}
func GenerateAccessToken(userId string, secretKey *config.Config) (string, error) {
	sKey := secretKey.JWTSecretKey

	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(sKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefresToken(userId string)(string, error){
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	stringToken := hex.EncodeToString(b)
	return stringToken, nil
}


 
	

