package minio

import (
	"context"
	"time"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMinio(config *env.Config) (*minio.Client, error) {
	// Initialize minio client object.
	useSSL := config.MINIO_USESSL == "true"
	minioClient, err := minio.New(config.MINIO_ENDPOINT, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MINIO_ACCESSKEY, config.MINIO_SECRETKEY, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if exists, err := minioClient.BucketExists(ctx, config.MINIO_BUCKET); err != nil {
		return nil, err
	} else if !exists {
		err = minioClient.MakeBucket(ctx, config.MINIO_BUCKET, minio.MakeBucketOptions{Region: config.MINIO_REGION})
		if err != nil {
			return nil, err
		}
	}
	// Set the client to the repository
	repository.DB.Min = minioClient
	return minioClient, nil
}
