package cloudstorageuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cloudstorageuc/cloudstoragedto"
	"github.com/localpaas/localpaas/services/aws/s3"
)

func (uc *CloudStorageUC) TestCloudStorageConn(
	ctx context.Context,
	auth *basedto.Auth,
	req *cloudstoragedto.TestCloudStorageConnReq,
) (*cloudstoragedto.TestCloudStorageConnResp, error) {
	setting, err := uc.SettingRepo.GetByID(ctx, uc.DB, nil, base.SettingTypeCloudProvider, req.Provider.ID, true)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	cloudConf, err := setting.AsCloudProvider()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	s3Client, err := s3.NewClient(ctx, &s3.Config{
		AccessKeyID:     cloudConf.AWS.AccessKeyID,
		SecretAccessKey: cloudConf.AWS.SecretKey.MustGetPlain(),
		Region:          gofn.Coalesce(req.S3.Region, cloudConf.AWS.Region),
		Endpoint:        req.S3.Endpoint,
		Bucket:          req.S3.Bucket,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	_, err = s3Client.HeadBucket(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cloudstoragedto.TestCloudStorageConnResp{}, nil
}
