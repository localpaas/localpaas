package gittokenuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
)

func (uc *GitTokenUC) UpdateGitToken(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.UpdateGitTokenReq,
) (*gittokendto.UpdateGitTokenResp, error) {
	req.Type = currentSettingType
	_, err := providers.UpdateSetting(ctx, uc.db, &req.UpdateSettingReq, &providers.UpdateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		PrepareUpdate: func(
			ctx context.Context,
			db database.Tx,
			data *providers.UpdateSettingData,
			pData *providers.PersistingSettingData,
		) error {
			pData.Setting.Name = gofn.Coalesce(req.Name, pData.Setting.Name)
			pData.Setting.Kind = gofn.Coalesce(string(req.Kind), pData.Setting.Kind)
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

	return &gittokendto.UpdateGitTokenResp{}, nil
}
