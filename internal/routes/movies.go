package routes

import (
	"log"

	"github.com/ctheil/pmdb-api/internal/controllers"
	"github.com/gin-gonic/gin"
)

func MovieRoutes(router *gin.RouterGroup) {
	mc, err := controllers.NewMovieController()
	if err != nil {
		log.Fatalf("could not init movie controller: %e", err)
	}
	r := router.Group("/movies")
	{
		r.GET("/trending", mc.GetTrendingMovies)
	}
}
