package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	APIKey   string `json:"api_key"`
	BasePath string `json:"base_path"`
	Images   struct {
		BaseUrl       string   `json:"base_url"`
		SecureBaseUrl string   `json:"secure_base_url"`
		BackdropSizes []string `json:"backdrop_sizes"`
		LogoSizes     []string `json:"logo_sizes"`
		PosterSizes   []string `json:"poster_sizes"`
		ProfileSizes  []string `json:"profile_sizes"`
	} `json:"images"`
}

func FetchConfig() (*Config, error) {
	config := Config{
		APIKey:   os.Getenv("tmdb_api_key"),
		BasePath: "https://api.themoviedb.org",
	}
	endpoint := "/3/configuration"
	req, err := http.NewRequest("GET", config.BasePath+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch config: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
