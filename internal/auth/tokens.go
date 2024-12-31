package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/ctheil/pmdb-api/internal/model"
	"github.com/ctheil/pmdb-api/internal/repository"
	"github.com/ctheil/pmdb-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type Token interface {
	IsValid() bool
}

type Claims interface{}

type AccessToken struct {
	Token  string
	Claims *services.UserClaims
}

type RefreshToken struct {
	Token  string
	User   model.User
	Claims *services.RefreshClaims
}

func ParseAccessToken(token string) (*AccessToken, error) {
	if token == "" {
		return &AccessToken{}, fmt.Errorf("no access token provided")
	}

	return &AccessToken{
		Token:  token,
		Claims: services.ParseAccessToken(token),
	}, nil
}

func ParseRefreshToken(token string, user model.User) (*RefreshToken, error) {
	if token == "" {
		return &RefreshToken{}, fmt.Errorf("no access token provided")
	}

	return &RefreshToken{
		Token:  token,
		Claims: services.ParseRefreshToken(token),
		User:   user,
	}, nil
}

func (t *AccessToken) IsValid() bool {
	return t.Claims.ExpiresAt.After(time.Now())
}

func (t *RefreshToken) IsValid() bool {
	return t.Claims.ExpiresAt.After(time.Now()) && t.Claims.Version == t.User.RefreshTokenVersion
}

type AuthService struct {
	TX *sqlx.Tx
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

func NewTokens(user model.User, user_repo *repository.UserRespoitory) (AuthTokens, error) {
	at := AuthTokens{}

	var err error
	at.AccessToken, err = NewAccessToken(user.Username, user.Id)
	if err != nil {
		return at, err
	}
	at.RefreshToken, err = NewRefreshToken(user, user_repo)
	if err != nil {
		return at, err
	}
	return at, nil
}

func NewAccessToken(username string, id uint) (string, error) {
	userClaims := services.UserClaims{
		Id:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)), // DEV
		},
	}

	signedAccessToken, err := services.NewAccessToken(userClaims)
	if err != nil {
		return "", err
	}
	return signedAccessToken, nil
}

func NewRefreshToken(user model.User, user_repo *repository.UserRespoitory) (string, error) {
	refreshClaims := services.RefreshClaims{
		Version: user.RefreshTokenVersion + 1,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)), // 30 days
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 2)), // 30 days // DEV 2 minutes
		},
	}

	signedRefreshToken, err := services.NewRefreshToken(refreshClaims)
	if err != nil {
		return "", err
	}
	user_repo.UpdateField(user, "refresh_token_version", user.RefreshTokenVersion+1)

	return signedRefreshToken, nil
}

func ClearTokenCookie(c *gin.Context, name string) {
	env := os.Getenv("ENV")
	prod := env == "Production"
	httpOnly := true
	secure := env == "Production"
	path := "/"
	domain := "https://pmdb.com"
	if !prod {
		domain = ""
	}

	c.SetCookie(name, "", 0, path, domain, secure, httpOnly)
}

func SetTokenCookies(c *gin.Context, accessToken, refreshToken string) {
	env := os.Getenv("ENV")
	prod := env == "Production"
	httpOnly := true
	secure := env == "Production"
	path := "/"
	domain := "https://pmdb.com"
	if !prod {
		domain = ""
	}

	// BUG: maxAge for refreshToken should not be this long...
	maxAge := int(time.Now().Year() * 10)
	if accessToken != "" {
		c.SetCookie("access_token", accessToken, maxAge, path, domain, secure, httpOnly)
	}

	if refreshToken != "" {
		c.SetCookie("refresh_token", refreshToken, maxAge, path, domain, secure, httpOnly)
	}
}

func GetTokenCookies(c *gin.Context) (access_token, refresh_token string, err error) {
	access_token, err = c.Cookie("access_token")
	if err != nil {
		return "", "", err
	}
	refresh_token, err = c.Cookie("refresh_token")
	return access_token, refresh_token, err
}
