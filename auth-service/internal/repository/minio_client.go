package repository

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/obadoraibu/go-auth/internal/config"
	"github.com/sirupsen/logrus"
)

type MinIOClient struct {
	Client *minio.Client
	Bucket string
}

func NewMinIOClient(cfg *config.MinIOConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		logrus.Error("Ошибка инициализации MinIO клиента: ", err)
		return nil, err
	}
	logrus.Info("MinIO клиент успешно инициализирован")

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		logrus.Error("Ошибка при проверке существования бакета: ", err)
		return nil, err
	}

	if !exists {
		logrus.Infof("Бакет %s не существует, пытаемся создать...", cfg.Bucket)
		err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			logrus.Error("Не удалось создать бакет: ", err)
			return nil, err
		}
		logrus.Infof("Бакет %s успешно создан", cfg.Bucket)
	} else {
		logrus.Infof("Бакет %s уже существует", cfg.Bucket)
	}

	return &MinIOClient{
		Client: client,
		Bucket: cfg.Bucket,
	}, nil
}

func (r *Repository) UploadScanToMinIO(ctx context.Context, objectName string, file multipart.File, size int64, contentType string) error {
	info, err := r.MinIO.Client.PutObject(
		ctx,
		r.MinIO.Bucket,
		objectName,
		file,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return nil
	}
	logrus.Infof("MRI uploaded to MinIO: %s (%d bytes)", info.Key, info.Size)
	return nil
}
