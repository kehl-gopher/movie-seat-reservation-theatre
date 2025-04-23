package minio

import (
	"bytes"
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

func UploadToMinio(min *minio.Client, filePath string, bucketName, contentType, objectName string, byt []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	reader := bytes.NewReader(byt)
	defer cancel()
	info, err := min.PutObject(ctx, bucketName, objectName, reader, int64(len(byt)), minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		return "", err
	}

	minioUrl := min.EndpointURL().Host + "/" + info.Bucket + "/" + objectName
	return minioUrl, nil
}

func GetObjectFromMinio(min *minio.Client) {

}
