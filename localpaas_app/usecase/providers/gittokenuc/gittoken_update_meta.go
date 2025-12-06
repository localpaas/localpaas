package gittokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
)

func (uc *GitTokenUC) UpdateGitTokenMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.UpdateGitTokenMetaReq,
) (*gittokendto.UpdateGitTokenMetaResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		tokenData := &updateGitTokenData{}
		err := uc.loadGitTokenDataForUpdateMeta(ctx, db, req, tokenData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		uc.prepareUpdatingGitTokenMeta(req, tokenData)
		return uc.persistGitTokenMeta(ctx, db, tokenData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gittokendto.UpdateGitTokenMetaResp{}, nil
}

func (uc *GitTokenUC) loadGitTokenDataForUpdateMeta(
	ctx context.Context,
	db database.IDB,
	req *gittokendto.UpdateGitTokenMetaReq,
	data *updateGitTokenData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeGitToken, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting

	return nil
}

func (uc *GitTokenUC) prepareUpdatingGitTokenMeta(
	req *gittokendto.UpdateGitTokenMetaReq,
	data *updateGitTokenData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	setting.UpdateVer++
	setting.UpdatedAt = timeNow

	if req.Status != nil {
		setting.Status = *req.Status
	}
	if req.ExpireAt != nil {
		setting.ExpireAt = *req.ExpireAt
	}
}

func (uc *GitTokenUC) persistGitTokenMeta(
	ctx context.Context,
	db database.IDB,
	data *updateGitTokenData,
) error {
	err := uc.settingRepo.Update(ctx, db, data.Setting,
		bunex.UpdateColumns("update_ver", "updated_at", "status", "expire_at"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
