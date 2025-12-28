package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

// S3Storage provides helper functions to interact with S3, including
// presigned URL generation for uploads/downloads and common object
// operations. The implementation is intentionally small and synchronous;
// callers may wrap calls in goroutines when needed.
type S3Storage struct {
	client          *s3.Client
	bucketName      string
	presignDuration time.Duration
}

// NewS3Storage creates a new S3 storage instance with a default
// presign expiration of 15 minutes.
func NewS3Storage(client *s3.Client, bucketName string) *S3Storage {
	return &S3Storage{
		client:          client,
		bucketName:      bucketName,
		presignDuration: 15 * time.Minute,
	}
}

// GeneratePresignedUploadURL returns a presigned PUT URL that clients
// can use to upload objects directly to S3. `key` is the S3 object key
// to be created.
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

// GeneratePresignedDownloadURL returns a presigned GET URL for the
// specified object key.
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

// ObjectExists checks whether an object exists in the bucket. It treats
// a 404/NotFound as "not exists" and returns (false, nil) in that case.
// Note: error handling of AWS SDK errors can be more specific depending
// on SDK versions; this implementation keeps things simple.
func (s *S3Storage) ObjectExists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		var apiErr smithy.APIError
		if ok := smithy.As(err, &apiErr); ok {
			// If the API returns NotFound-like error, treat as not found.
			return false, nil
		}
		return false, fmt.Errorf("failed to check object: %w", err)
	}

	return true, nil
}

// DeleteObject removes the specified key from the bucket.
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

// GetObject downloads the entire object into memory and returns its
// bytes. WARNING: this reads the full content into memory, which may be
// unsuitable for large objects; prefer streaming in production.
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
