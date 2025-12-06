package ssluc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/ssluc/ssldto"
)

func (uc *SslUC) UpdateSslMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *ssldto.UpdateSslMetaReq,
) (*ssldto.UpdateSslMetaResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		sslData := &updateSslData{}
		err := uc.loadSslDataForUpdateMeta(ctx, db, req, sslData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		uc.prepareUpdatingSslMeta(req, sslData)
		return uc.persistSslMeta(ctx, db, sslData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &ssldto.UpdateSslMetaResp{}, nil
}

func (uc *SslUC) loadSslDataForUpdateMeta(
	ctx context.Context,
	db database.IDB,
	req *ssldto.UpdateSslMetaReq,
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

	return nil
}

func (uc *SslUC) prepareUpdatingSslMeta(
	req *ssldto.UpdateSslMetaReq,
	data *updateSslData,
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

func (uc *SslUC) persistSslMeta(
	ctx context.Context,
	db database.IDB,
	data *updateSslData,
) error {
	err := uc.settingRepo.Update(ctx, db, data.Setting,
		bunex.UpdateColumns("update_ver", "updated_at", "status", "expire_at"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
