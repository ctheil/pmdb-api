package controllers

import (
	"fmt"
	"net/http"

	"github.com/ctheil/pmdb-api/internal/services"
	"github.com/gin-gonic/gin"
)

func GetTrendingMovies(c *gin.Context) {
	movies, err := services.FetchTrendingMovies()
	if err != nil {
		fmt.Printf("error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting trending movies."})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success!",
		"results": movies,
	})
}
