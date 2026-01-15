package media

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/chai2010/webp"
)

// R2Service handles Cloudflare R2 operations (S3-compatible)
type R2Service struct {
	client  *s3.S3
	bucket  string
	baseURL string
}

// NewR2Service creates a new Cloudflare R2 service instance
func NewR2Service() (*R2Service, error) {
	endpoint := os.Getenv("R2_ENDPOINT")
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucket := os.Getenv("R2_BUCKET")
	baseURL := os.Getenv("IMAGE_BASE_URL")

	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
		return nil, fmt.Errorf("missing R2 configuration: R2_ENDPOINT, R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_BUCKET required")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("auto"),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create R2 session: %v", err)
	}

	return &R2Service{
		client:  s3.New(sess),
		bucket:  bucket,
		baseURL: baseURL,
	}, nil
}

// NewS3Service creates a new R2 service (alias for backwards compatibility)
func NewS3Service() (*R2Service, error) {
	return NewR2Service()
}

// S3Service is an alias for R2Service (backwards compatibility)
type S3Service = R2Service

// UploadImage uploads an image to R2 and returns the public URL
func (s *R2Service) UploadImage(img image.Image, folder string) (string, error) {
	// Generate filename with timestamp
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("alerty_%d.webp", timestamp)
	key := filepath.Join(folder, filename)

	// Convert image to WebP format in memory
	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Quality: 80}); err != nil {
		return "", fmt.Errorf("failed to encode image to webp: %v", err)
	}

	// Upload to R2
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket:       aws.String(s.bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(buf.Bytes()),
		ContentType:  aws.String("image/webp"),
		CacheControl: aws.String("max-age=31536000"), // 1 year cache
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %v", err)
	}

	// Return public URL via CDN
	url := fmt.Sprintf("%s/%s", s.baseURL, key)
	return url, nil
}

// UploadRawFile uploads a raw file to R2 (for temporary processing)
func (s *R2Service) UploadRawFile(data []byte, folder, filename string) (string, error) {
	key := filepath.Join(folder, filename)

	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/octet-stream"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload raw file to R2: %v", err)
	}

	url := fmt.Sprintf("%s/%s", s.baseURL, key)
	return url, nil
}

// DeleteFile deletes a file from R2
func (s *R2Service) DeleteFile(key string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
