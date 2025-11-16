package ssluc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

func (uc *SslUC) DeleteSsl(
	ctx context.Context,
	auth *basedto.Auth,
	req *ssldto.DeleteSslReq,
) (*ssldto.DeleteSslResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		registryAuthData := &deleteSslData{}
		err := uc.loadSslDataForDelete(ctx, db, req, registryAuthData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingSslData{}
		uc.prepareDeletingSsl(registryAuthData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &ssldto.DeleteSslResp{}, nil
}

type deleteSslData struct {
	Setting *entity.Setting
}

func (uc *SslUC) loadSslDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *ssldto.DeleteSslReq,
	data *deleteSslData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, req.ID,
		bunex.SelectFor("UPDATE OF setting"),
		bunex.SelectWhere("setting.type = ?", base.SettingTypeSsl),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *SslUC) prepareDeletingSsl(
	data *deleteSslData,
	persistingData *persistingSslData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
