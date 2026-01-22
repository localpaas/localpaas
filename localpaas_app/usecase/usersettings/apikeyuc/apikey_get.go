package apikeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

func (uc *APIKeyUC) GetAPIKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.GetAPIKeyReq,
) (*apikeydto.GetAPIKeyResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, currentSettingType, req.ID, false,
		bunex.SelectWhere("setting.deleted_at IS NULL"),
		bunex.SelectWhere("setting.object_id = ?", auth.User.ID),
		bunex.SelectRelation("ObjectUser", bunex.SelectWithDeleted()),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := apikeydto.TransformAPIKey(setting, auth.User.ID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.GetAPIKeyResp{
		Data: resp,
	}, nil
}
