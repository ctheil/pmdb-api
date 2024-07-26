package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func BuildReq(url, method string, reader io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	authStr := fmt.Sprintf("Bearer %s", os.Getenv("tmdb_api_key"))
	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Authorization", authStr)

	return req, nil
}

func ExecReq(req *http.Request) (*http.Response, error) {
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error client.Do(req): %e", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to execute request: %s", resp.Status)
	}

	return resp, nil
}

func ReqToJSON(data io.ReadCloser, out any) error {
	decoder := json.NewDecoder(data)

	if err := decoder.Decode(&out); err != nil {
		fmt.Printf("[ReqToJSON]: error unmarshaling...\n")
		return err
	}

	return nil
}
