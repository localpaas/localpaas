package githubappuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/githubappuc/githubappdto"
)

func (uc *GithubAppUC) UpdateGithubAppMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *githubappdto.UpdateGithubAppMetaReq,
) (*githubappdto.UpdateGithubAppMetaResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		appData := &updateGithubAppData{}
		err := uc.loadGithubAppDataForUpdateMeta(ctx, db, req, appData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		uc.prepareUpdatingGithubAppMeta(req, appData)
		return uc.persistGithubAppMeta(ctx, db, appData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &githubappdto.UpdateGithubAppMetaResp{}, nil
}

func (uc *GithubAppUC) loadGithubAppDataForUpdateMeta(
	ctx context.Context,
	db database.IDB,
	req *githubappdto.UpdateGithubAppMetaReq,
	data *updateGithubAppData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeGithubApp, req.ID, false,
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

func (uc *GithubAppUC) prepareUpdatingGithubAppMeta(
	req *githubappdto.UpdateGithubAppMetaReq,
	data *updateGithubAppData,
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

func (uc *GithubAppUC) persistGithubAppMeta(
	ctx context.Context,
	db database.IDB,
	data *updateGithubAppData,
) error {
	err := uc.settingRepo.Update(ctx, db, data.Setting,
		bunex.UpdateColumns("update_ver", "updated_at", "status", "expire_at"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
