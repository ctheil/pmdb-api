package routes

import (
	"database/sql"

	"github.com/ctheil/pmdb-api/internal/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, db *sql.DB) {
	uc := controllers.NewUserController(db)
	r := router.Group("/auth")
	{
		r.GET("/user/:id", uc.GetUser)
		r.POST("/user", uc.InsertUser)
	}
}
