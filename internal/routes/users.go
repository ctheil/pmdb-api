package routes

import (
	"github.com/ctheil/pmdb-api/internal/controllers"
	"github.com/ctheil/pmdb-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func UserRoutes(router *gin.RouterGroup, tx *sqlx.Tx) {
	uc := controllers.NewUserController(tx)
	r := router.Group("/user")
	{
		r.GET("/", middleware.Authenticate(tx), uc.GetUser)
	}
}
