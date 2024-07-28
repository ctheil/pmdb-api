package routes

import (
	"github.com/ctheil/pmdb-api/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func UserRoutes(router *gin.RouterGroup, tx *sqlx.Tx) {
	uc := controllers.NewUserController(tx)
	r := router.Group("/auth")
	{
		r.GET("/user/:id", uc.GetUser)
		r.POST("/user", uc.InsertUser)
		r.POST("/login", uc.Login)
	}
}
