package e2e

import (
	"github.com/bdragon300/tusgo"
	"github.com/gocastsian/roham/filer/adapter/tusdclient"
	"github.com/stretchr/testify/assert"
	"net/url"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {

	//Create a temporary test file
	testFileContent := []byte("This is a test file")
	filePath := "test.txt"
	err := os.WriteFile(filePath, testFileContent, 0644)
	assert.NoError(t, err)
	defer os.Remove(filePath)

	baseURL, _ := url.Parse("http://localhost:5006/uploads/")
	cl := tusdclient.New(baseURL)

	f, err := os.Open("/home/nimamleo/Downloads/ostan.zip")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	u := tusdclient.CreateUploadFromFile(f, cl)

	stream := tusgo.NewUploadStream(cl, u)
	if err = tusdclient.UploadWithRetry(stream, f); err != nil {
		panic(err)
	}

	// example of uploadID that returned by minio
	// The first part, 4bbe6710877de4dc20b1b227ab68adc4, is typically the MinIO object key.
	// 4bbe6710877de4dc20b1b227ab68adc4+NzAxZDA4MjctZmRjMy00YzdlLTkxOTQtMDJlYTQxZmYzZDgyLjc1ZGY0ZTFlLTFhNWEtNDA1Ni04Y2QxLTkzNDA1MzkxOWJkOXgxNzQ1MzY4Nzg4MDc5MTYxNzU3
}
