package storage

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"field-service/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client     *minio.Client
	bucketName string
	publicURL  string
}

func NewMinIOClient(cfg config.MinIOConfig) (IStorageClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio new client: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("bucket exists: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("make bucket: %w", err)
		}
	}
	// policy public di-set tiap startup (idempotent) biar bucket selalu kebaca dari luar
	policy := fmt.Sprintf(
		`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`,
		cfg.BucketName,
	)
	if err := client.SetBucketPolicy(ctx, cfg.BucketName, policy); err != nil {
		return nil, fmt.Errorf("set bucket policy: %w", err)
	}

	return &MinIOClient{client: client, bucketName: cfg.BucketName, publicURL: cfg.PublicURL}, nil
}

func (m *MinIOClient) UploadFile(ctx context.Context, filename string, data []byte) (string, error) {
	_, err := m.client.PutObject(ctx, m.bucketName, filename, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: http.DetectContentType(data),
	})
	if err != nil {
		return "", fmt.Errorf("put object: %w", err)
	}
	return fmt.Sprintf("%s/%s/%s", m.publicURL, m.bucketName, filename), nil
}
