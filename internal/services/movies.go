package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ctheil/pmdb-api/internal/config"
)

type Movie struct {
	OriginalTitle    string   `json:"original_title"`
	Overview         string   `json:"overview"`
	Genres           []string `json:"genres"`
	OriginalLanguage string   `json:"original_language"`
	PosterPath       string   `json:"poster_path"`
	Title            string   `json:"title"`
	Status           string   `json:"status"`
	IMDB_ID          string   `json:"imdb_id"`
	Id               int      `json:"id"`
	ReleaseDate      string   `json:"release_date"`
	Runtime          int      `json:"runtime"`
}

type TMDBMoviesResponse struct {
	Page    int     `json:"page"`
	Results []Movie `json:"results"`
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
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading response body")
		return nil, err

	}

	tmdbMovies := TMDBMoviesResponse{}
	if err := json.Unmarshal(body, &tmdbMovies); err != nil {
		fmt.Println("Error unmarshaling data... %e", err)
		return nil, err
	}

	return tmdbMovies.Results, nil
}
