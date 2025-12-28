package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

// S3Storage handles S3 operations
type S3Storage struct {
	client       *s3.Client
	bucketName   string
	presignDuration time.Duration
}

// NewS3Storage creates a new S3 storage instance
func NewS3Storage(client *s3.Client, bucketName string) *S3Storage {
	return &S3Storage{
		client:          client,
		bucketName:      bucketName,
		presignDuration: 15 * time.Minute,
	}
}

// GeneratePresignedUploadURL generates a presigned URL for uploading
func (s *S3Storage) GeneratePresignedUploadURL(ctx context.Context, key string) (string, error) {
	presignClient := s3.NewPresignFromClient(s.client)

	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(s.presignDuration.Seconds()) * time.Second
	})

	if err != nil {
		return "", fmt.Errorf("failed to presign URL: %w", err)
	}

	return request.URL, nil
}

// GeneratePresignedDownloadURL generates a presigned URL for downloading
func (s *S3Storage) GeneratePresignedDownloadURL(ctx context.Context, key string) (string, error) {
	presignClient := s3.NewPresignFromClient(s.client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(s.presignDuration.Seconds()) * time.Second
	})

	if err != nil {
		return "", fmt.Errorf("failed to presign URL: %w", err)
	}

	return request.URL, nil
}

// ObjectExists checks if an object exists in S3
func (s *S3Storage) ObjectExists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		var apiErr smithy.APIError
		if awsErr := new(smithy.APIError); err != awsErr {
			// Check if 404
			return false, nil
		}
		return false, fmt.Errorf("failed to check object: %w", err)
	}

	return true, nil
}

// DeleteObject deletes an object from S3
func (s *S3Storage) DeleteObject(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// GetObject retrieves an object from S3
func (s *S3Storage) GetObject(ctx context.Context, key string) ([]byte, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer result.Body.Close()

	// Read body (simplified - in production, stream to avoid loading entire file into memory)
	buf := make([]byte, *result.ContentLength)
	_, err = result.Body.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return buf, nil
}
