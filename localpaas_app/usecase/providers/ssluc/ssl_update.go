package ssluc

import (
	"context"
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/ssluc/ssldto"
)

func (uc *SslUC) UpdateSsl(
	ctx context.Context,
	auth *basedto.Auth,
	req *ssldto.UpdateSslReq,
) (*ssldto.UpdateSslResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		sslData := &updateSslData{}
		err := uc.loadSslDataForUpdate(ctx, db, req, sslData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingSslData{}
		uc.prepareUpdatingSsl(req.SslBaseReq, sslData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &ssldto.UpdateSslResp{}, nil
}

type updateSslData struct {
	Setting *entity.Setting
}

func (uc *SslUC) loadSslDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *ssldto.UpdateSslReq,
	data *updateSslData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeSsl, req.ID, false,
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
		conflictSetting, _ := uc.settingRepo.GetByName(ctx, db, base.SettingTypeSsl, req.Name, false)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist("SSL").
				WithMsgLog("ssl '%s' already exists", conflictSetting.Name)
		}
	}

	return nil
}

func (uc *SslUC) prepareUpdatingSsl(
	req *ssldto.SslBaseReq,
	data *updateSslData,
	persistingData *persistingSslData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	setting.UpdateVer++
	setting.UpdatedAt = timeNow
	setting.Name = gofn.Coalesce(req.Name, setting.Name)

	ssl := &entity.Ssl{
		Certificate: req.Certificate,
		PrivateKey:  entity.NewEncryptedField(req.PrivateKey),
		KeySize:     req.KeySize,
		Provider:    req.Provider,
		Email:       req.Email,
		Expiration:  req.Expiration,
	}
	setting.MustSetData(ssl)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
