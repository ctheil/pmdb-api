package config

import (
	"encoding/json"
	"io"
)

func ReqToJSON(data io.ReadCloser, out any) error {
	body, err := io.ReadAll(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &out); err != nil {
		return err
	}

	return nil
}
