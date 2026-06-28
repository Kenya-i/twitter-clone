package repository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Kenya-i/twitter-clone/internal/config"
	"github.com/Kenya-i/twitter-clone/internal/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Storage struct {
	client   *s3.Client
	bucket   string
	endpoint string
}

func NewS3Storage(cfg *config.Config) (domain.ImageStorage, error) {
	client := s3.New(s3.Options{
		Region:       "us-east-1",
		BaseEndpoint: aws.String(cfg.S3Endpoint),
		UsePathStyle: true,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.S3AccessKey, cfg.S3SecretKey, "",
		),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(cfg.S3Bucket),
	})
	if err != nil {
		var alreadyOwned *s3.BucketAlreadyOwnedByYou
		var alreadyExists *s3.BucketAlreadyExists
		if !errors.As(err, &alreadyOwned) && !errors.As(err, &alreadyExists) {
			return nil, err
		}
	}

	return &s3Storage{client: client, bucket: cfg.S3Bucket, endpoint: cfg.S3Endpoint}, nil
}

func (s *s3Storage) Upload(file io.Reader, filename string, contentType string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, filename), nil
}
