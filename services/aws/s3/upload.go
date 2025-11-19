package s3

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/h2non/filetype"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	defaultContentType = "application/octet-stream" // Default content type for unknown files
	uploadPartSize     = 10 * 1024 * 1024           // 10MB
)

func (client *Client) Upload(ctx context.Context, bucketName string, objectKey string, data []byte) error {
	kind, err := filetype.Match(data)
	if err != nil {
		return apperrors.Wrap(err)
	}
	contentType := defaultContentType
	if kind != filetype.Unknown {
		contentType = kind.MIME.Value
	}

	uploader := manager.NewUploader(client.client, func(u *manager.Uploader) {
		// Let S3 split the data into parts and upload them in parallel
		if len(data) > uploadPartSize {
			u.PartSize = uploadPartSize
		}
	})
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
