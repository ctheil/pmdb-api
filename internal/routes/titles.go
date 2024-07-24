package routes

import (
	"log"

	"github.com/ctheil/pmdb-api/internal/controllers"
	"github.com/gin-gonic/gin"
)

func TitleRoutes(router *gin.RouterGroup) {
	tc, err := controllers.NewTitleController()
	if err != nil {
		log.Fatalf("could not init movie controller: %e", err)
	}
	r := router.Group("/titles")
	{
		r.GET("/trending", tc.GetTrendingTitles)
		r.GET("/details/:id", tc.GetDetailsById)
	}
}
