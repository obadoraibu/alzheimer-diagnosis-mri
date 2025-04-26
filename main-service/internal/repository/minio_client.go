// repository/minio.go
package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/obadoraibu/go-auth/internal/config"
	"github.com/sirupsen/logrus"
)

type MinIOClient struct {
	internal *minio.Client // для Put/Get/List и т.-д. внутри кластера
	public   *minio.Client // только для генерации ссылок
	bucket   string
}

// создаём два клиента с одинаковыми cred’ами, но разными endpoint’ами
func NewMinIOClient(cfg *config.MinIOConfig) (*MinIOClient, error) {
	creds := credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, "")

	internal, err := minio.New(cfg.InternalEndpoint, &minio.Options{
		Creds:  creds,
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	public, err := minio.New(cfg.PublicEndpoint, &minio.Options{
		Creds:        creds,
		Secure:       cfg.UseSSL,
		Region:       "us-east-1", // чтобы SDK не делал лишний запрос
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := internal.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := internal.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
		logrus.Infof("bucket %s created", cfg.Bucket)
	}

	return &MinIOClient{
		internal: internal,
		public:   public,
		bucket:   cfg.Bucket,
	}, nil
}

func (r *Repository) UploadScanToMinIO(
	ctx context.Context,
	objectName string,
	file multipart.File,
	size int64,
	contentType string,
) error {
	_, err := r.MinIO.internal.PutObject(ctx, r.MinIO.bucket, objectName, file, size,
		minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (r *Repository) PresignedGetObject(objectName string) (*url.URL, error) {
	if objectName == "" {
		return nil, fmt.Errorf("empty object name")
	}
	return r.MinIO.public.PresignedGetObject(
		context.Background(),
		r.MinIO.bucket,
		objectName,
		15*time.Minute,
		nil,
	)
}
