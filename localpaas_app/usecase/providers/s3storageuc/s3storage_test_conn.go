package s3storageuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/s3storageuc/s3storagedto"
	"github.com/localpaas/localpaas/services/aws/s3"
)

func (uc *S3StorageUC) TestS3StorageConn(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.TestS3StorageConnReq,
) (*s3storagedto.TestS3StorageConnResp, error) {
	s3Client, err := s3.NewClient(ctx, &s3.Config{
		AccessKeyID:     req.AccessKeyID,
		SecretAccessKey: req.SecretKey,
		Endpoint:        req.Endpoint,
		Region:          req.Region,
		Bucket:          req.Bucket,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	_, err = s3Client.HeadBucket(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.TestS3StorageConnResp{}, nil
}
