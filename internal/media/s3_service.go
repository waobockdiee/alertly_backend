package media

import (
	"bytes"
	"fmt"
	"image"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/chai2010/webp"
)

const (
	S3BucketName = "alertly-images-production"
	S3Region     = "us-west-2"
)

// S3Service handles S3 operations
type S3Service struct {
	client *s3.S3
	bucket string
}

// NewS3Service creates a new S3 service instance
func NewS3Service() (*S3Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(S3Region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return &S3Service{
		client: s3.New(sess),
		bucket: S3BucketName,
	}, nil
}

// UploadImage uploads an image to S3 and returns the public URL
func (s *S3Service) UploadImage(img image.Image, folder string) (string, error) {
	// Generate filename with timestamp
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("alerty_%d.webp", timestamp)
	key := filepath.Join(folder, filename)

	// Convert image to WebP format in memory
	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Quality: 80}); err != nil {
		return "", fmt.Errorf("failed to encode image to webp: %v", err)
	}

	// Upload to S3
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket:       aws.String(s.bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(buf.Bytes()),
		ContentType:  aws.String("image/webp"),
		CacheControl: aws.String("max-age=31536000"), // 1 year cache
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	// Return public URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, S3Region, key)
	return url, nil
}

// UploadRawFile uploads a raw file to S3 (for temporary processing)
func (s *S3Service) UploadRawFile(data []byte, folder, filename string) (string, error) {
	key := filepath.Join(folder, filename)

	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/octet-stream"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload raw file to S3: %v", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, S3Region, key)
	return url, nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(key string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
