package routes

import (
	"github.com/ctheil/pmdb-api/internal/controllers"
	"github.com/gin-gonic/gin"
)

func MovieRoutes(router *gin.RouterGroup) {
	r := router.Group("/movies")
	{
		r.GET("/trending", controllers.GetTrendingMovies)
	}
}
