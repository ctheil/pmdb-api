package routes

import (
	"github.com/ctheil/pmdb-api/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AuthRoutes(router *gin.RouterGroup, tx *sqlx.Tx) {
	ac := controllers.NewAuthController(tx)
	r := router.Group("/auth")
	{
		r.POST("/login", ac.Login)
		// r.POST("/signup", middleware.AuthenticateTokens(tx), ac.Signup) // PROTECT
		r.POST("/signup", ac.Signup)
	}
}
