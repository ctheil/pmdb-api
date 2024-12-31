package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ctheil/pmdb-api/internal/auth"
	"github.com/ctheil/pmdb-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type OAuthController struct {
	TX    *sqlx.Tx
	OAuth *services.OAuth
}

func NewOAuthController(tx *sqlx.Tx) *OAuthController {
	oauth, err := services.NewOAuth()
	if err != nil {
		log.Fatalf("Error generating oauth config... %e", err)
	}

	return &OAuthController{TX: tx, OAuth: oauth}
}

func (a *OAuthController) GetAuthUrl(c *gin.Context) {
	oauth, err := services.NewOAuth()
	if err != nil {
		fmt.Printf("Error generating oauth config: %s", err)

		log.Fatalf("Failed to generate oauth config: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	url := fmt.Sprintf("%s?%s", oauth.CFG.Auth_URL, oauth.GetOAuthParams())
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved oauth url", "url": url})
}

type OAuthUserData struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (a *OAuthController) GetAuthToken(c *gin.Context) {
	code := c.Request.URL.Query().Get("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Authorization code must be provided"})
		return
	}
	// TODO: implement retries via context?
	oauth_resp, err := a.OAuth.FetchAuth()
	if err != nil {
		fmt.Printf("Error fetching oauth: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not verify authentication via OAuth provider."})
		return
	}
	userData := services.ParseOAuthUserToken(oauth_resp.ID_Token)
	newToken, err := services.NewOAuthUserToken(userData.UserData, a.OAuth.CFG.Token_Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating new oauth token."})
		return
	}

	auth.SetTokenCookies(c, newToken, "")
	/*
	* TODO:
	*   Store User in DB
	*   Use ACCESS/REFRESH token system
	* */

	c.JSON(http.StatusOK, gin.H{"message": "Successfully found user.", "user": userData})
}

func (a *OAuthController) GetLoggedIn(c *gin.Context) {
	newToken, err := a.OAuth.VerifyToken(c)
	if err != nil {
		log.Print("Either missing token or no user data in token.")
		c.JSON(http.StatusOK, gin.H{"message": "No user data found in token, or no token provided", "loggedIn": false})
		return
	}
	auth.SetTokenCookies(c, newToken, "")
	c.JSON(http.StatusOK, gin.H{"loggedIn": true})
}

func (a *OAuthController) PostLogout(c *gin.Context) {
	auth.ClearTokenCookie(c, "access_token")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}