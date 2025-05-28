package tusdadapter

import (
	"errors"
	"github.com/gocastsian/roham/filer/service/filestorage"
	"github.com/gocastsian/roham/types"
	"github.com/tus/tusd/pkg/handler"
	"strconv"
	"strings"
)

type CompleteUploadsHookEvent handler.HookEvent

func (e *CompleteUploadsHookEvent) ConvertToCreateFileMetadataInput() (filestorage.CreateFileMetadataInput, error) {

	storageID, err := strconv.ParseInt(e.Upload.MetaData["TARGET-STORAGE-ID"], 10, 64)
	if err != nil {
		return filestorage.CreateFileMetadataInput{}, err
	}

	// 4bbe6710877de4dc20b1b227ab68adc4+NzAxZDA4MjctZmRjMy00YzdlLTkxOTQtMDJlYTQxZmYzZDgyLjc1ZGY0ZTFlLTFhNWEtNDA1Ni04Y2QxLTkzNDA1MzkxOWJkOXgxNzQ1MzY4Nzg4MDc5MTYxNzU3

	parts := strings.Split(e.Upload.ID, "+")

	if len(parts) == 0 {
		return filestorage.CreateFileMetadataInput{}, errors.New("invalid upload id")
	}

	fileKey := parts[0]

	return filestorage.CreateFileMetadataInput{
		TargetStorageID: types.ID(storageID),
		FileKey:         fileKey,
		FileName:        e.Upload.MetaData["filename"],
		MimeType:        e.Upload.MetaData["filetype"],
		Size:            e.Upload.Size,
	}, nil
}
