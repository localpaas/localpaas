package awss3uc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awss3uc/awss3dto"
	"github.com/localpaas/localpaas/services/aws/s3"
)

func (uc *AWSS3UC) TestAWSS3Conn(
	ctx context.Context,
	auth *basedto.Auth,
	req *awss3dto.TestAWSS3ConnReq,
) (*awss3dto.TestAWSS3ConnResp, error) {
	awsSetting, err := uc.settingRepo.GetByID(ctx, uc.db, base.SettingTypeAWS, req.Cred.ID, true)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	awsConf, err := awsSetting.AsAWS()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	s3Client, err := s3.NewClient(ctx, &s3.Config{
		AccessKeyID:     awsConf.AccessKeyID,
		SecretAccessKey: awsConf.SecretKey.MustGetPlain(),
		Region:          gofn.Coalesce(req.Region, awsConf.Region),
		Endpoint:        req.Endpoint,
		Bucket:          req.Bucket,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	_, err = s3Client.HeadBucket(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &awss3dto.TestAWSS3ConnResp{}, nil
}
