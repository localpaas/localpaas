package s3storagedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type GetS3StorageReq struct {
	ID string `json:"-"`
}

func NewGetS3StorageReq() *GetS3StorageReq {
	return &GetS3StorageReq{}
}

func (req *GetS3StorageReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetS3StorageResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *S3StorageResp    `json:"data"`
}

type S3StorageResp struct {
	ID              string                        `json:"id"`
	Name            string                        `json:"name"`
	AccessKeyID     string                        `json:"accessKeyId"`
	Region          string                        `json:"region"`
	Bucket          string                        `json:"bucket"`
	ProjectAccesses []*S3StorageProjectAccessResp `json:"projectAccesses"`
}

type S3StorageBaseResp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AccessKeyID string `json:"accessKeyId"`
	Region      string `json:"region"`
	Bucket      string `json:"bucket"`
}

type S3StorageProjectAccessResp struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	Allowed     bool                      `json:"allowed"`
	AppAccesses []*S3StorageAppAccessResp `json:"appAccesses"`
}

type S3StorageAppAccessResp struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Allowed bool   `json:"allowed"`
}

func TransformS3Storage(s3Storage *entity.S3Storage) (resp *S3StorageResp, err error) {
	if err = copier.Copy(&resp, &s3Storage); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.ProjectAccesses, err = TransformS3StorageObjectAccesses(s3Storage.ObjectAccesses)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

func TransformS3StorageObjectAccesses(accesses []*entity.ACLPermission) (
	resp []*S3StorageProjectAccessResp, err error) {
	mapProjectResp := make(map[string]*S3StorageProjectAccessResp)
	for _, access := range accesses {
		if access.SubjectType != base.SubjectTypeProject || access.SubjectProject == nil {
			continue
		}
		projResp := &S3StorageProjectAccessResp{
			ID:      access.SubjectID,
			Name:    access.SubjectProject.Name,
			Allowed: access.Actions.Read || access.Actions.Write || access.Actions.Delete,
		}
		resp = append(resp, projResp)
		mapProjectResp[access.SubjectID] = projResp
	}
	for _, access := range accesses {
		if access.SubjectType != base.SubjectTypeApp || access.SubjectApp == nil {
			continue
		}
		projResp := mapProjectResp[access.SubjectApp.ProjectID]
		if projResp != nil {
			projResp.AppAccesses = append(projResp.AppAccesses, &S3StorageAppAccessResp{
				ID:      access.SubjectID,
				Name:    access.SubjectApp.Name,
				Allowed: access.Actions.Read || access.Actions.Write || access.Actions.Delete,
			})
		}
	}
	return resp, nil
}

func TransformS3StorageBase(s3Storage *entity.S3Storage) (resp *S3StorageBaseResp, err error) {
	if err = copier.Copy(&resp, &s3Storage); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

func TransformS3StoragesBase(s3Storages []*entity.S3Storage) ([]*S3StorageBaseResp, error) {
	return basedto.TransformObjectSlice(s3Storages, TransformS3StorageBase) //nolint:wrapcheck
}
