package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type Client struct {
	client *s3.Client
	config *Config
}

type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	Region          string
	Bucket          string
}

func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithDefaultRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(opts *s3.Options) {
		if cfg.Endpoint != "" {
			opts.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	return &Client{client: s3Client, config: cfg}, nil
}

func (client *Client) HeadBucket(ctx context.Context) (*s3.HeadBucketOutput, error) {
	result, err := client.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(client.config.Bucket),
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return result, nil
}
