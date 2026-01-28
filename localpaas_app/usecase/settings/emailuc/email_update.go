package emailuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/emailuc/emaildto"
)

func (uc *EmailUC) UpdateEmail(
	ctx context.Context,
	auth *basedto.Auth,
	req *emaildto.UpdateEmailReq,
) (*emaildto.UpdateEmailResp, error) {
	req.Type = currentSettingType
	_, err := settings.UpdateSetting(ctx, uc.db, &req.UpdateSettingReq, &settings.UpdateSettingData{
		SettingRepo:       uc.settingRepo,
		VerifyingName:     req.Name,
		DefaultMustUnique: true,
		PrepareUpdate: func(
			ctx context.Context,
			db database.Tx,
			data *settings.UpdateSettingData,
			pData *settings.PersistingSettingData,
		) error {
			pData.Setting.Name = gofn.Coalesce(req.Name, pData.Setting.Name)
			email := &entity.Email{}
			if req.SMTP != nil {
				pData.Setting.Kind = string(base.EmailKindSMTP)
				email.SMTP = &entity.SMTPConf{
					Host:        req.SMTP.Host,
					Port:        req.SMTP.Port,
					Username:    req.SMTP.Username,
					DisplayName: req.SMTP.DisplayName,
					Password:    entity.NewEncryptedField(req.SMTP.Password),
					SSL:         req.SMTP.SSL,
				}
			}
			if req.HTTP != nil {
				pData.Setting.Kind = string(base.EmailKindHTTP)
				email.HTTP = &entity.HTTPMailConf{
					Endpoint:    req.HTTP.Endpoint,
					Method:      req.HTTP.Method,
					ContentType: req.HTTP.ContentType,
					Headers:     req.HTTP.Headers,
					BodyMapping: req.HTTP.BodyMapping,
				}
			}
			err := pData.Setting.SetData(email)
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &emaildto.UpdateEmailResp{}, nil
}
