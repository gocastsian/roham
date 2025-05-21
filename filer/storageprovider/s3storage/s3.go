package s3storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gocastsian/roham/filer/storageprovider"
	"io"
	"time"
)

type Storage struct {
	s3  *s3.S3
	cfg storageprovider.StorageConfig
}

func New(cfg storageprovider.StorageConfig) (*Storage, error) {

	awsConfig := aws.NewConfig().
		WithRegion("ir").
		WithEndpoint(cfg.Endpoint).
		WithCredentials(credentials.NewStaticCredentials(
			cfg.AccessKey,
			cfg.SecretKey,
			cfg.Region, // Leave empty unless using STS/OpenID
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

func (s *Storage) Config() storageprovider.StorageConfig {
	return s.cfg
}
