package services

import (
	"fmt"
	"net/http"

	"github.com/ctheil/pmdb-api/internal/config"
	"github.com/ctheil/pmdb-api/internal/model"
)

type TMDBMoviesResponse struct {
	Page    int           `json:"page"`
	Results []model.Title `json:"results"`
}

func FetchTrendingMovies() ([]Movie, error) {
	apiUrl := "https://api.themoviedb.org/3/trending/movie/day?language=en-US"

	req, err := config.TmdbRequest(apiUrl, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("error formatting request: ", err)
	}

	c := &http.Client{}

	res, err := c.Do(req)
	if err != nil {
		fmt.Println("Error executing request...")
		return nil, err
	} else if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("tmdb returned status: %s", res.Status)
	}
	defer res.Body.Close()

	tmdbMovies := TMDBMoviesResponse{}
	if err := config.ReqToJSON(res.Body, tmdbMovies); err != nil {
		return nil, err
	}

	return tmdbMovies.Results, nil
}
