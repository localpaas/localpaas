package basicauthuc

import (
	"context"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc/basicauthdto"
)

func (uc *BasicAuthUC) UpdateBasicAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *basicauthdto.UpdateBasicAuthReq,
) (*basicauthdto.UpdateBasicAuthResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		basicAuthData := &updateBasicAuthData{}
		err := uc.loadBasicAuthDataForUpdate(ctx, db, req, basicAuthData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingBasicAuthData{}
		uc.prepareUpdatingBasicAuth(req.BasicAuthBaseReq, basicAuthData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &basicauthdto.UpdateBasicAuthResp{}, nil
}

type updateBasicAuthData struct {
	Setting *entity.Setting
}

func (uc *BasicAuthUC) loadBasicAuthDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *basicauthdto.UpdateBasicAuthReq,
	data *updateBasicAuthData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeBasicAuth, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting

	// If name changes, validate the new one
	if req.Name != "" && !strings.EqualFold(setting.Name, req.Name) {
		conflictSetting, _ := uc.settingRepo.GetByName(ctx, db, base.SettingTypeBasicAuth, req.Name, false)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist("BasicAuth").
				WithMsgLog("basic auth '%s' already exists", conflictSetting.Name)
		}
	}

	return nil
}

func (uc *BasicAuthUC) prepareUpdatingBasicAuth(
	req *basicauthdto.BasicAuthBaseReq,
	data *updateBasicAuthData,
	persistingData *persistingBasicAuthData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	if req.Name != "" {
		setting.Name = req.Name
	}

	basicAuth := &entity.BasicAuth{
		Username: req.Username,
		Password: req.Password,
	}
	setting.MustSetData(basicAuth.MustEncrypt())

	setting.UpdatedAt = timeNow
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
