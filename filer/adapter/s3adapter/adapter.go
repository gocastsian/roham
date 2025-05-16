package s3adapter

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"time"
)

type Adapter struct {
	s3 *s3.S3
}

func NewSession(cfg *aws.Config) (*session.Session, error) {
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func New(cfg *aws.Config) (*Adapter, error) {

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		s3: s3.New(sess),
	}, err
}

func (a *Adapter) S3() *s3.S3 {
	return a.s3
}

func (a *Adapter) GetFileContent(ctx context.Context, bucketName, key string) (io.ReadCloser, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	result, err := a.s3.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return result.Body, nil
}

func (a *Adapter) GeneratePreSignedURL(bucketName, key string, duration time.Duration) (string, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	req, _ := a.s3.GetObjectRequest(input)
	urlStr, err := req.Presign(duration)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}

	return urlStr, nil
}

func (a *Adapter) MakeStorage(ctx context.Context, name string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	_, err := a.s3.CreateBucketWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", name, err)
	}

	return nil
}
