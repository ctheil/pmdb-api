package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ctheil/pmdb-api/internal/model"
	"github.com/ctheil/pmdb-api/internal/services"
	"github.com/gin-gonic/gin"
)

/*
* NOTE: SETUP
* Types, Init, Setup
* */
type TitleController struct {
	// Config  *config.Config
	Service *services.TitleService
}

func NewTitleController() (*TitleController, error) {
	s, err := services.NewTitleService()
	if err != nil {
		log.Fatalf("failed to init new title service: %e", err)
		return nil, err
	}

	return &TitleController{Service: s}, nil
}

func (tc *TitleController) GetTrendingTitles(c *gin.Context) {
	include_images := c.Query("include_images")
	image_size := c.Query("image_size")
	domain, err := tc.Service.GetDomain(c.Query("domain"))
	if err != nil {
		fmt.Printf("error: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	region := "en-US"

	endpoint := fmt.Sprintf("/3/trending/%s/day?language=%s", domain, region)
	titleResp := model.TitleResponse{}
	if err := tc.Service.TMDBRequest(endpoint, "GET", nil, &titleResp); err != nil {
		fmt.Printf("error executing TMDBRequest: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if include_images == "true" {
		s := 0
		size, err := strconv.Atoi(image_size)
		if err != nil {
			fmt.Println("error converting image_size to int, fallback to 0")
		}
		tc.Service.BuildPosterPaths(titleResp.Results, size|s)
		tc.Service.BuildBackdropPaths(titleResp.Results, size|s)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success!",
		"results": titleResp.Results,
	})
}

func (tc *TitleController) GetDetailsById(c *gin.Context) {
	id := c.Param("id")
	include := c.Query("include") // NOTE: v = video c = credit p = providers
	includes := tc.Service.ExtractIncludes(include)
	region := "en-US"

	details_endpt := "/3/movie/" + id + "?append_to_response=" + includes + "&language=" + region

	details := model.TitleDetails{}
	if err := tc.Service.TMDBRequest(details_endpt, "GET", nil, &details); err != nil {
		log.Printf("error in TMDBRequest: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}
	tc.Service.BuiltProfilePaths(details.Credits.Cast, 3)
	tc.Service.BuildProviderLogoPaths(details.WatchProviders, 3)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success!",
		"results": details,
	})
}
