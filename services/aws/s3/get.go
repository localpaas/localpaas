package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (client *Client) GetObject(ctx context.Context, bucketName string, objectKey string) ([]byte, error) {
	result, err := client.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return data, nil
}
