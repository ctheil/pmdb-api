package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func TmdbRequest(url string, method string, reader io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("error formatting request: %e", err)
	}

	authStr := fmt.Sprintf("Bearer %s", os.Getenv("tmdb_api_key"))

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Authorization", authStr)

	return req, nil
}
