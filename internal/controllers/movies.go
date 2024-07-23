package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/ctheil/pmdb-api/internal/config"
	"github.com/ctheil/pmdb-api/internal/model"
	"github.com/gin-gonic/gin"
)

type MovieController struct {
	Config *config.Config
}

func NewMovieController() (*MovieController, error) {
	cfg, err := config.FetchConfig()
	if err != nil {
		return nil, err
	}

	return &MovieController{Config: cfg}, nil
}

func (mc *MovieController) BuildReq(path, method string, reader io.Reader) (*http.Request, error) {
	url := mc.Config.BasePath + path

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	authStr := fmt.Sprintf("Bearer %s", os.Getenv("tmdb_api_key"))
	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Authorization", authStr)

	return req, nil
}

func (mc *MovieController) ExecReq(req *http.Request) (*http.Response, error) {
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to execute request: %s", resp.Status)
	}

	return resp, nil
}

func (mc *MovieController) GetTrendingMovies(c *gin.Context) {
	include_images := c.Query("include_images")
	image_size := c.Query("image_size")
	req, err := mc.BuildReq("/3/trending/movie/day?language=en-US", "GET", nil)
	if err != nil {
		fmt.Printf("error: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting trending movies."})
		return
	}
	resp, err := mc.ExecReq(req)
	if err != nil {
		fmt.Printf("error: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting trending movies."})
		return
	}

	titleResp := model.TitleResponse{}
	fmt.Printf("response: %v\n", resp.Status)
	if err := config.ReqToJSON(resp.Body, &titleResp); err != nil {
		fmt.Printf("error: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting trending movies."})
		return
	}

	if include_images == "true" {
		s := 0
		size, err := strconv.Atoi(image_size)
		if err != nil {
			fmt.Println("error converting image_size to int, fallback to 0")
		}
		mc.buildPosterPath(titleResp.Results, size|s)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success!",
		"results": titleResp.Results,
	})
}

func (mc *MovieController) buildPosterPath(titles []model.Title, size int) {
	len := len(mc.Config.Images.PosterSizes)
	if size > len {
		size = len - 1
	}
	for i, title := range titles {
		titles[i].PosterPath = mc.Config.Images.BaseUrl + mc.Config.Images.PosterSizes[size] + title.PosterPath
	}
}
