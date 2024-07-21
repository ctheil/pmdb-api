package main

import (
	"log"

	"github.com/ctheil/pmdb-api/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file... %e", err)
	}

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		routes.MovieRoutes(v1)
	}

	r.Run(":8080")
}
