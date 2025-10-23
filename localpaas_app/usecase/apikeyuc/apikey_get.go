package apikeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/apikeyuc/apikeydto"
)

func (uc *APIKeyUC) GetAPIKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.GetAPIKeyReq,
) (*apikeydto.GetAPIKeyResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, req.ID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	apiKey, err := setting.ParseAPIKey()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	userMap := make(map[string]*entity.User)
	if apiKey != nil {
		user, err := uc.userRepo.GetByID(ctx, uc.db, apiKey.ActingUser.ID, bunex.SelectWithDeleted())
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		userMap[user.ID] = user
	}

	resp, err := apikeydto.TransformAPIKey(setting, userMap)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.GetAPIKeyResp{
		Data: resp,
	}, nil
}
