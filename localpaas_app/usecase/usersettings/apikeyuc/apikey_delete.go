package apikeyuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

func (uc *APIKeyUC) DeleteAPIKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.DeleteAPIKeyReq,
) (*apikeydto.DeleteAPIKeyResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		apiKeyData := &deleteAPIKeyData{}
		err := uc.loadAPIKeyDataForDelete(ctx, db, auth, req, apiKeyData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingAPIKeyData{}
		uc.prepareDeletingAPIKey(apiKeyData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.DeleteAPIKeyResp{}, nil
}

type deleteAPIKeyData struct {
	Setting *entity.Setting
}

func (uc *APIKeyUC) loadAPIKeyDataForDelete(
	ctx context.Context,
	db database.IDB,
	auth *basedto.Auth,
	req *apikeydto.DeleteAPIKeyReq,
	data *deleteAPIKeyData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeAPIKey, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
		bunex.SelectWhere("setting.deleted_at IS NULL"),
		bunex.SelectWhere("setting.object_id = ?", auth.User.ID),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *APIKeyUC) prepareDeletingAPIKey(
	data *deleteAPIKeyData,
	persistingData *persistingAPIKeyData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
