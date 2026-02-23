package imagebuilduc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuilduc/imagebuilddto"
)

const (
	currentSettingType    = base.SettingTypeImageBuild
	currentSettingVersion = entity.CurrentImageBuildVersion

	defaultName = "image build settings"
)

func (uc *ImageBuildUC) CreateImageBuild(
	ctx context.Context,
	auth *basedto.Auth,
	req *imagebuilddto.CreateImageBuildReq,
) (*imagebuilddto.CreateImageBuildResp, error) {
	req.Type = currentSettingType
	resp, err := uc.CreateSetting(ctx, &req.CreateSettingReq, &settings.CreateSettingData{
		Version: currentSettingVersion,
		PrepareCreation: func(
			ctx context.Context,
			db database.Tx,
			data *settings.CreateSettingData,
			pData *settings.PersistingSettingCreationData,
		) error {
			pData.Setting.Name = defaultName
			err := pData.Setting.SetData(req.ToEntity())
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
		AfterPersisting: func(
			ctx context.Context,
			db database.Tx,
			data *settings.CreateSettingData,
			pData *settings.PersistingSettingCreationData,
		) error {
			return uc.ensureSettingIsUniqueInScope(ctx, db, &req.BaseSettingReq)
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &imagebuilddto.CreateImageBuildResp{
		Data: resp.Data,
	}, nil
}

func (uc *ImageBuildUC) ensureSettingIsUniqueInScope(
	ctx context.Context,
	db database.IDB,
	req *settings.BaseSettingReq,
) (err error) {
	switch req.Scope { //nolint:exhaustive
	case base.SettingScopeGlobal:
		err = uc.SettingRepo.EnsureUniqueGlobally(ctx, db, req.Type)
	case base.SettingScopeProject:
		err = uc.SettingRepo.EnsureUniqueInProject(ctx, db, req.Type, req.ObjectID)
	default:
		return apperrors.NewUnsupported().WithMsgLog("setting scope '%v' is not supported", req.Scope)
	}
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
