package fetcher

import (
	"bytes"
	"fmt"
	"github.com/xiaoshouchen/openapi-generator/internal/model"
	"io"
	"net/http"
)

type Apifox struct {
	config model.Fetcher
}

func NewApifoxFetcher(config model.Fetcher) *Apifox {
	return &Apifox{config: config}
}

func (f *Apifox) Bytes() ([]byte, error) {
	return GetData(f.config.Apifox.Token, f.config.Apifox.ProjectId)
}

func GetData(token string, projectId int64) ([]byte, error) {
	url := fmt.Sprintf("https://api.apifox.com/v1/projects/%d/export-openapi", projectId)

	payload := []byte(`{
        "scope": {
            "type": "ALL","excludedByTags": []
        },
        "options": {
            "includeApifoxExtensionProperties": false,
            "addFoldersToTags": false
        },
        "oasVersion": "3.0",
        "exportFormat": "JSON"
    }`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-Apifox-Api-Version", "2024-03-28")
	req.Header.Set("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
