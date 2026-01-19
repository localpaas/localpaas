package gittokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
)

const (
	currentSettingType    = base.SettingTypeGitToken
	currentSettingVersion = entity.CurrentGitTokenVersion
)

func (uc *GitTokenUC) CreateGitToken(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.CreateGitTokenReq,
) (*gittokendto.CreateGitTokenResp, error) {
	req.Type = currentSettingType
	resp, err := providers.CreateSetting(ctx, uc.db, &req.CreateSettingReq, &providers.CreateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		Version:       currentSettingVersion,
		PrepareCreation: func(ctx context.Context, db database.Tx, data *providers.CreateSettingData,
			pData *providers.PersistingSettingCreationData) error {
			pData.Setting.Kind = string(req.Kind)
			pData.Setting.ExpireAt = req.ExpireAt
			err := pData.Setting.SetData(&entity.GitToken{
				User:    req.User,
				Token:   entity.NewEncryptedField(req.Token),
				BaseURL: req.BaseURL,
			})
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gittokendto.CreateGitTokenResp{
		Data: resp.Data,
	}, nil
}
