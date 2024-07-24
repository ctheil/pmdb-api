package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ctheil/pmdb-api/internal/config"
	"github.com/ctheil/pmdb-api/internal/model"
	"github.com/gin-gonic/gin"
)

/*
* NOTE: SETUP
* Types, Init, Setup
* */
type TitleController struct {
	Config *config.Config
}

func NewTitleController() (*TitleController, error) {
	cfg, err := config.FetchConfig()
	if err != nil {
		return nil, err
	}

	return &TitleController{Config: cfg}, nil
}

/*
* NOTE: UTILS
* builtReq, execReq, buildPosterPaths
* */
func (tc *TitleController) buildReq(path, method string, reader io.Reader) (*http.Request, error) {
	url := tc.Config.BasePath + path

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	authStr := fmt.Sprintf("Bearer %s", os.Getenv("tmdb_api_key"))
	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Authorization", authStr)

	return req, nil
}

func (tc *TitleController) execReq(req *http.Request) (*http.Response, error) {
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

func (tc *TitleController) buildPosterPaths(titles []model.Title, size int) {
	len := len(tc.Config.Images.PosterSizes)
	if size > len {
		size = len - 1
	}
	for i, title := range titles {
		titles[i].PosterPath = tc.Config.Images.BaseUrl + tc.Config.Images.PosterSizes[size] + title.PosterPath
	}
}

func (tc *TitleController) buildBackdropPaths(titles []model.Title, size int) {
	len := len(tc.Config.Images.BackdropSizes)
	if size > len {
		size = len - 1
	}
	for i, title := range titles {
		path := tc.Config.Images.BaseUrl + tc.Config.Images.BackdropSizes[size] + title.BackdropPath
		titles[i].BackdropPath = path
	}
}

func (tc *TitleController) buildLogoPaths(companies []model.ProductionCompany, size int) {
	len := len(tc.Config.Images.LogoSizes)
	if size > len {
		size = len - 1
	}
	for i, c := range companies {
		if c.LogoPath == "" {
			continue
		}
		path := tc.Config.Images.BaseUrl + tc.Config.Images.LogoSizes[size] + c.LogoPath
		companies[i].LogoPath = path
	}
}

func (tc *TitleController) getDomain(domain string) (string, error) {
	switch domain {
	case "movie":
		return domain, nil
	case "tv":
		return domain, nil
	case "person":
		return domain, nil
	default:
		return "", fmt.Errorf("invalid domain: %s", domain)
	}
}

/*
* NOTE: CONTROLLER FUNCTIONS
* GetTrendingMovies
* */

func (tc *TitleController) TMDBRequest(endpoint, method string, reader io.Reader, out interface{}) error {
	req, err := tc.buildReq(endpoint, method, reader)
	if err != nil {
		return fmt.Errorf("error building request: %e", err)
	}
	resp, err := tc.execReq(req)
	if err != nil {
		return fmt.Errorf("error executing request: %e", err)
	}

	if err := config.ReqToJSON(resp.Body, out); err != nil {
		return fmt.Errorf("error decoding body: %e", err)
	}
	return nil
}

func (tc *TitleController) GetTrendingTitles(c *gin.Context) {
	include_images := c.Query("include_images")
	image_size := c.Query("image_size")
	domain, err := tc.getDomain(c.Query("domain"))
	if err != nil {
		fmt.Printf("error: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	region := "en-US"

	endpoint := fmt.Sprintf("/3/trending/%s/day?language=%s", domain, region)
	titleResp := model.TitleResponse{}
	tc.TMDBRequest(endpoint, "GET", nil, &titleResp)
	if include_images == "true" {
		s := 0
		size, err := strconv.Atoi(image_size)
		if err != nil {
			fmt.Println("error converting image_size to int, fallback to 0")
		}
		tc.buildPosterPaths(titleResp.Results, size|s)
		tc.buildBackdropPaths(titleResp.Results, size|s)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success!",
		"results": titleResp.Results,
	})
}

func (tc *TitleController) GetDetailsById(c *gin.Context) {
	id := c.Param("id")

	// NOTE: Only allows access to en-US titles. Should store the user's reguion in the session for easy access to incluse in the request in the future.
	region := "en-US"
	details_endpt := fmt.Sprintf("/3/movie/%s?language=%s", id, region)
	details := model.TitleDetails{}
	if err := tc.TMDBRequest(details_endpt, "GET", nil, &details); err != nil {
		log.Printf("error in TMDBRequest: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	video_endpt := fmt.Sprintf("/3/movie/%s/videos?language=%s", id, region)
	videos := model.Videos{}
	if err := tc.TMDBRequest(video_endpt, "GET", nil, &videos); err != nil {
		log.Printf("error in TMDBRequest: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}
	tc.buildLogoPaths(details.ProductionCompanies, 4)

	details.Videos = videos

	c.JSON(http.StatusOK, gin.H{
		"message": "Success!",
		"results": details,
	})
}
