package fetcher

import "github.com/xiaoshouchen/openapi-go-generator/internal/model"

type Fetcher interface {
	Bytes() ([]byte, error)
}

func NewFetcher(fetcherConfig model.Fetcher) (Fetcher, error) {
	var fetcher Fetcher
	switch fetcherConfig.From {
	case "file":
		fetcher = NewFileFetcher(fetcherConfig.File)
	case "apifox":
		fetcher = NewApifoxFetcher(fetcherConfig)
	}
	return fetcher, nil
}
