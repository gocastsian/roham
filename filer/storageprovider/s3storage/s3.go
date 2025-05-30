package s3storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gocastsian/roham/filer/storageprovider"
)

type Storage struct {
	s3  *s3.S3
	cfg storageprovider.StorageConfig
}

func New(cfg storageprovider.StorageConfig) (*Storage, error) {

	awsConfig := aws.NewConfig().
		WithRegion(cfg.Region).
		WithEndpoint(cfg.Endpoint).
		WithCredentials(credentials.NewStaticCredentials(
			cfg.AccessKey,
			cfg.SecretKey,
			"", // Leave empty unless using STS/OpenID
		)).
		WithS3ForcePathStyle(true).
		WithDisableSSL(true)

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	return &Storage{
		s3:  s3.New(sess),
		cfg: cfg,
	}, err
}

func (s *Storage) S3() *s3.S3 {
	return s.s3
}

func (s *Storage) GetFile(ctx context.Context, storageName, fileKey string) (io.ReadCloser, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(storageName),
		Key:    aws.String(fileKey),
	}

	result, err := s.s3.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return result.Body, nil
}

func (s *Storage) GeneratePreSignedURL(storageName, key string, duration time.Duration) (string, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(storageName),
		Key:    aws.String(key),
	}

	req, _ := s.s3.GetObjectRequest(input)
	urlStr, err := req.Presign(duration)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}

	return urlStr, nil
}

func (s *Storage) MakeStorage(ctx context.Context, name string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	_, err := s.s3.CreateBucketWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", name, err)
	}

	return nil
}

func (s *Storage) MoveFileToStorage(fileKey, fromStorageName, toStorageName string) error {
	// If source and destination are the same, no need to move
	if fromStorageName == toStorageName {
		return nil
	}

	// Move both the main file and its .info file
	filesToMove := []string{fileKey, fileKey + ".info"}

	for _, file := range filesToMove {
		// Format the copy source correctly (without leading slash)
		copySource := fmt.Sprintf("%s/%s", fromStorageName, file)

		// Copy the object to the destination bucket
		_, err := s.s3.CopyObject(&s3.CopyObjectInput{
			Bucket:     aws.String(toStorageName),
			CopySource: aws.String(copySource),
			Key:        aws.String(file),
		})
		if err != nil {
			// If copy fails, try to clean up any previously copied files
			for _, cleanupFile := range filesToMove {
				if cleanupFile != file {
					_, _ = s.s3.DeleteObject(&s3.DeleteObjectInput{
						Bucket: aws.String(toStorageName),
						Key:    aws.String(cleanupFile),
					})
				}
			}
			return fmt.Errorf("failed to copy file %s from %s to %s: %w", file, fromStorageName, toStorageName, err)
		}

		// Wait until the object is copied
		err = s.s3.WaitUntilObjectExists(&s3.HeadObjectInput{
			Bucket: aws.String(toStorageName),
			Key:    aws.String(file),
		})
		if err != nil {
			// If waiting fails, try to clean up all copied objects
			for _, cleanupFile := range filesToMove {
				_, _ = s.s3.DeleteObject(&s3.DeleteObjectInput{
					Bucket: aws.String(toStorageName),
					Key:    aws.String(cleanupFile),
				})
			}
			return fmt.Errorf("waiting for copied object %s failed: %w", file, err)
		}

		// Delete the object from the source bucket
		_, err = s.s3.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(fromStorageName),
			Key:    aws.String(file),
		})
		if err != nil {
			// If deletion fails, we should log this but not return an error
			// since the file was successfully copied
			return fmt.Errorf("failed to delete original file %s from %s: %w", file, fromStorageName, err)
		}
	}

	return nil
}

func (s *Storage) Config() storageprovider.StorageConfig {
	return s.cfg
}
