package gittokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
)

func (uc *GitTokenUC) DeleteGitToken(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.DeleteGitTokenReq,
) (*gittokendto.DeleteGitTokenResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		tokenData := &deleteGitTokenData{}
		err := uc.loadGitTokenDataForDelete(ctx, db, req, tokenData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingGitTokenData{}
		uc.prepareDeletingGitToken(tokenData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gittokendto.DeleteGitTokenResp{}, nil
}

type deleteGitTokenData struct {
	Setting *entity.Setting
}

func (uc *GitTokenUC) loadGitTokenDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *gittokendto.DeleteGitTokenReq,
	data *deleteGitTokenData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeGitToken, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *GitTokenUC) prepareDeletingGitToken(
	data *deleteGitTokenData,
	persistingData *persistingGitTokenData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
