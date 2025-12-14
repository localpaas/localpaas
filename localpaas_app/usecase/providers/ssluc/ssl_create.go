package ssluc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/ssluc/ssldto"
)

func (uc *SslUC) CreateSsl(
	ctx context.Context,
	auth *basedto.Auth,
	req *ssldto.CreateSslReq,
) (*ssldto.CreateSslResp, error) {
	sslData := &createSslData{}
	err := uc.loadSslData(ctx, uc.db, req, sslData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingSslData{}
	uc.preparePersistingSsl(req.SslBaseReq, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &ssldto.CreateSslResp{
		Data: &basedto.ObjectIDResp{ID: createdItem.ID},
	}, nil
}

type createSslData struct {
}

func (uc *SslUC) loadSslData(
	ctx context.Context,
	db database.IDB,
	req *ssldto.CreateSslReq,
	_ *createSslData,
) error {
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeSsl, req.Name, false)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("Ssl").
			WithMsgLog("ssl setting '%s' already exists", req.Name)
	}

	return nil
}

type persistingSslData struct {
	settingservice.PersistingSettingData
}

func (uc *SslUC) preparePersistingSsl(
	req *ssldto.SslBaseReq,
	persistingData *persistingSslData,
) {
	timeNow := timeutil.NowUTC()
	dbSsl := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeSsl,
		Status:    base.SettingStatusActive,
		Name:      req.Name,
		Kind:      string(req.Provider),
		Version:   entity.CurrentSslVersion,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	ssl := &entity.Ssl{
		Certificate: req.Certificate,
		PrivateKey:  entity.NewEncryptedField(req.PrivateKey),
		KeySize:     req.KeySize,
		Provider:    req.Provider,
		Email:       req.Email,
	}
	dbSsl.MustSetData(ssl)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, dbSsl)
}

func (uc *SslUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingSslData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
