package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ctheil/pmdb-api/internal/model"
	"github.com/ctheil/pmdb-api/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrServerError  = errors.New("internal server error")
)

type AuthError struct {
	Type    error
	Message string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("%v: %s", e.Type, e.Message)
}

func AuthenticateRequest(c *gin.Context, tx *sqlx.Tx) *AuthError {
	access_token, refresh_token, err := GetTokenCookies(c)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return &AuthError{
				Type:    ErrUnauthorized,
				Message: "no auth tokens in header.",
			}
		}
		return &AuthError{
			Type:    ErrServerError,
			Message: "error parsing cookies from header.",
		}
	}
	accessToken, err := ParseAccessToken(access_token)
	if err != nil {
		return &AuthError{
			Type:    ErrUnauthorized,
			Message: "failed to parse access token or no token.",
		}
	}

	user := model.User{
		Username: accessToken.Claims.Username,
		Id:       accessToken.Claims.Id,
	}
	// generate and send new access token
	new_access_token, err := NewAccessToken(user.Username, user.Id)
	if err != nil {
		return &AuthError{
			Type:    ErrServerError,
			Message: "failed to generate new access token",
		}
	}
	if accessToken.IsValid() {
		c.Set("user", user)
		SetTokenCookies(c, new_access_token, "")
		return nil
	}

	// else check refreshToken, if valid(), new tokens, c.Next(),
	user_repo := repository.NewUserRepository(tx)
	user, err = user_repo.GetByUsername(accessToken.Claims.Username)
	if err != nil {
		return &AuthError{
			Type:    ErrUnauthorized,
			Message: "user not found",
		}
		// c.JSON(http.StatusNotFound, gin.H{"message": "Unauthorized: could not find user"})
		// return
	}
	refreshToken, err := ParseRefreshToken(refresh_token, user)
	if err != nil {
		return &AuthError{
			Type:    ErrUnauthorized,
			Message: "failed to parase refresh token",
		}
	}

	if !refreshToken.IsValid() {
		return &AuthError{
			Type:    ErrUnauthorized,
			Message: "Please login",
		}
	}

	// Refresh token is valid, refresh both tokens and continue;
	new_tokens, err := NewTokens(user, user_repo)
	if err != nil {
		return &AuthError{
			Type:    ErrServerError,
			Message: "failed to generate new tokens",
		}
	}

	c.JSON(http.StatusOK, gin.H{"refresh_token": new_tokens.RefreshToken, "access_token": new_tokens.AccessToken})

	c.Set("user", user)
	c.Next()
	return nil
}
