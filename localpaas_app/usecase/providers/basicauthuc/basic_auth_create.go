package basicauthuc

import (
	"context"
	"errors"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc/basicauthdto"
)

func (uc *BasicAuthUC) CreateBasicAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *basicauthdto.CreateBasicAuthReq,
) (*basicauthdto.CreateBasicAuthResp, error) {
	basicAuthData := &createBasicAuthData{}
	err := uc.loadBasicAuthData(ctx, uc.db, req, basicAuthData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingBasicAuthData{}
	uc.preparePersistingBasicAuth(req.BasicAuthBaseReq, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &basicauthdto.CreateBasicAuthResp{
		Data: &basedto.ObjectIDResp{ID: createdItem.ID},
	}, nil
}

type createBasicAuthData struct {
}

func (uc *BasicAuthUC) loadBasicAuthData(
	ctx context.Context,
	db database.IDB,
	req *basicauthdto.CreateBasicAuthReq,
	_ *createBasicAuthData,
) error {
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeBasicAuth, req.Name, false)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("BasicAuth").
			WithMsgLog("basic auth setting '%s' already exists", req.Name)
	}

	return nil
}

type persistingBasicAuthData struct {
	settingservice.PersistingSettingData
}

func (uc *BasicAuthUC) preparePersistingBasicAuth(
	req *basicauthdto.BasicAuthBaseReq,
	persistingData *persistingBasicAuthData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeBasicAuth,
		Status:    base.SettingStatusActive,
		Name:      req.Name,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	basicAuth := &entity.BasicAuth{
		Username: req.Username,
		Password: entity.NewEncryptedField(req.Password),
	}
	setting.MustSetData(basicAuth)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}

func (uc *BasicAuthUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingBasicAuthData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
