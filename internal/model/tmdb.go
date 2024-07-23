package model

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ctheil/pmdb-api/internal/config"
)

type TMDB struct {
	Config struct {
		Images struct {
			BaseUrl       string   `json:"base_url"`
			SecureBaseUrl string   `json:"secure_base_url"`
			BackdropSizes []string `json:"backdrop_sizes"`
			LogoSizes     []string `json:"logo_sizes"`
			PosterSizes   []string `json:"poster_sizes"`
		} `json:"images"`
	}
}

func (t *TMDB) TmdbRequest(url string, method string, reader io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("error formatting request: %e", err)
	}

	authStr := fmt.Sprintf("Bearer %s", os.Getenv("tmdb_api_key"))

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Authorization", authStr)

	return req, nil
}

func (t *TMDB) Setup() error {
	url := "https://api.themoviedb.org/3/configuration"
	req, err := t.TmdbRequest(url, "GET", nil)
	if err != nil {
		log.Fatalf("error getting config from tmbd: %e", err)
	}

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return err
	} else if res.StatusCode < 200 || res.StatusCode > 299 {
		return err
	}

	config.ReqToJSON(res.Body, t.Config)
	return nil
}
