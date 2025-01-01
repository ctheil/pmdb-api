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

//	type OAuthUserData struct {
//		Email   string `json:"email"`
//		Name    string `json:"name"`
//		Picture string `json:"picture"`
//	}
type OAuthUserClaims struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret string
}

func NewTokenService() *TokenService {
	return &TokenService{os.Getenv("TOKEN_SECRET")}
}

func (t *TokenService) NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(t.secret))
}

type RefreshClaims struct {
	Version int `json:"version"`
	jwt.RegisteredClaims
}

func (t *TokenService) NewRefreshToken(claims RefreshClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(t.secret))
}

func (t *TokenService) ParseAccessToken(accessToken string) *UserClaims {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secret), nil
	})

	return parsedAccessToken.Claims.(*UserClaims)
}

func (t *TokenService) ParseRefreshToken(refreshToken string) *RefreshClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secret), nil
	})

	return parsedRefreshToken.Claims.(*RefreshClaims)
}

func (t *TokenService) ParseOAuthUserToken(oauthToken string) *OAuthUserClaims {
	parsedOauthToken, _ := jwt.ParseWithClaims(oauthToken, &OAuthUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secret), nil
	})
	return parsedOauthToken.Claims.(*OAuthUserClaims)
}

func (t *TokenService) NewOAuthUserToken(currentClaims *OAuthUserClaims) (string, error) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, currentClaims)

	return newToken.SignedString([]byte(t.secret))
}
