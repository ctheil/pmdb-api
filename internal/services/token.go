package services

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type OAuthUserData struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
type OAuthUserClaims struct {
	UserData OAuthUserData
	jwt.RegisteredClaims
}

func NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

type RefreshClaims struct {
	Version int `json:"version"`
	jwt.RegisteredClaims
}

func NewRefreshToken(claims RefreshClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAccessToken(accessToken string) *UserClaims {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedAccessToken.Claims.(*UserClaims)
}

func ParseRefreshToken(refreshToken string) *RefreshClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedRefreshToken.Claims.(*RefreshClaims)
}

func ParseOAuthUserToken(oauthToken string) *OAuthUserClaims {
	parsedOauthToken, _ := jwt.ParseWithClaims(oauthToken, &OAuthUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	return parsedOauthToken.Claims.(*OAuthUserClaims)
}

func NewOAuthUserToken(user OAuthUserData, secret string) (string, error) {
	claims := OAuthUserClaims{UserData: user}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return newToken.SignedString([]byte(secret))
}
