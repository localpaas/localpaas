package apikeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

func (uc *APIKeyUC) ListAPIKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.ListAPIKeyReq,
) (*apikeydto.ListAPIKeyResp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
		ExtraLoadOpts: []bunex.SelectQueryOption{
			bunex.SelectWhere("setting.deleted_at IS NULL"),
			bunex.SelectWhere("setting.object_id = ?", req.ObjectID),
			bunex.SelectRelation("ObjectUser", bunex.SelectWithDeleted()),
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := apikeydto.TransformAPIKeys(resp.Data, req.ObjectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.ListAPIKeyResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}
