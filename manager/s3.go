package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dblooman/centralis/resource"
)

// S3Storage implements the Storage interface using AWS S3
type S3Storage struct {
	bucket string
	svc    *s3.Client
}

// NewS3Storage creates a new S3Storage instance
func NewS3Storage(ctx context.Context, bucket, region string) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	svc := s3.NewFromConfig(cfg)

	return &S3Storage{
		bucket: bucket,
		svc:    svc,
	}, nil
}

// Save saves the resource data to S3
func (s *S3Storage) Save(ctx context.Context, data resource.ResourceData) error {
	key := filepath.Join(data.Type, data.ID+".json")
	data.CreatedAt = time.Now()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = s.svc.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(jsonData),
	})
	return err
}

// Load loads the resource data from S3
func (s *S3Storage) Load(ctx context.Context, resourceType, id string) (resource.ResourceData, error) {
	key := filepath.Join(resourceType, id+".json")
	result, err := s.svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return resource.ResourceData{}, err
	}

	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return resource.ResourceData{}, err
	}

	var data resource.ResourceData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return resource.ResourceData{}, err
	}

	return data, nil
}

// Delete deletes the resource data from S3
func (s *S3Storage) Delete(ctx context.Context, resourceType, id string) error {
	key := filepath.Join(resourceType, id+".json")
	_, err := s.svc.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

// List lists all resource IDs of a given type
func (s *S3Storage) List(ctx context.Context, resourceType string) ([]string, error) {
	prefix := resourceType + "/"
	result, err := s.svc.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(result.Contents))
	for _, item := range result.Contents {
		id := path.Base(*item.Key)
		ids = append(ids, id[:len(id)-len(".json")])
	}

	return ids, nil
}
