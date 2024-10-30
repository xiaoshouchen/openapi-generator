package fetcher

import "os"

type FileFetcher struct {
	filePath string
}

func NewFileFetcher(path string) *FileFetcher {
	return &FileFetcher{
		filePath: path,
	}
}

func (f *FileFetcher) Bytes() ([]byte, error) {
	return os.ReadFile(f.filePath)
}
