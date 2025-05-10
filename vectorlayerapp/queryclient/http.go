package queryclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type QueryClient struct {
}

func New() QueryClient {
	return QueryClient{}
}

func (q QueryClient) DownloadShapeFile(fileKey string) ([]byte, error) {
	baseUrl := "http://127.0.0.1:5005"
	encodedKey := url.PathEscape(fileKey)

	fullUrl := fmt.Sprintf("%s/v1/files/%s/download", baseUrl, encodedKey)

	resp, err := http.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}
