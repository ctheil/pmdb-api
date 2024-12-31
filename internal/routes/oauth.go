package routes

import (
	"github.com/ctheil/pmdb-api/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func OAuthRoutes(router *gin.RouterGroup, tx *sqlx.Tx) {
	ac := controllers.NewOAuthController(tx)
	r := router.Group("/oauth")
	{
		r.GET("/url", ac.GetAuthUrl)
		r.GET("/token", ac.GetAuthToken)
		r.GET("/logged_in", ac.GetLoggedIn)
		r.POST("/logout", ac.PostLogout)
		// r.GET("/auth/token")
	}
}
