package githubappuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/githubappuc/githubappdto"
)

func (uc *GithubAppUC) DeleteGithubApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *githubappdto.DeleteGithubAppReq,
) (*githubappdto.DeleteGithubAppResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		appData := &deleteGithubAppData{}
		err := uc.loadGithubAppDataForDelete(ctx, db, req, appData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingGithubAppData{}
		uc.prepareDeletingGithubApp(appData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &githubappdto.DeleteGithubAppResp{}, nil
}

type deleteGithubAppData struct {
	Setting *entity.Setting
}

func (uc *GithubAppUC) loadGithubAppDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *githubappdto.DeleteGithubAppReq,
	data *deleteGithubAppData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeGithubApp, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *GithubAppUC) prepareDeletingGithubApp(
	data *deleteGithubAppData,
	persistingData *persistingGithubAppData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
