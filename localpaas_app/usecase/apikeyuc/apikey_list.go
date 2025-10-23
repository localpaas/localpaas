package apikeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/entityutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/apikeyuc/apikeydto"
)

func (uc *APIKeyUC) ListAPIKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.ListAPIKeyReq,
) (*apikeydto.ListAPIKeyResp, error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.type = ?", base.SettingTypeAPIKey),
	}

	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhereGroup(
				bunex.SelectWhere("setting.name ILIKE ?", keyword),
			),
		)
	}

	settings, paging, err := uc.settingRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	userIDs := make([]string, 0, len(settings))
	for _, setting := range settings {
		apiKey, err := setting.ParseAPIKey()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		if apiKey != nil {
			userIDs = append(userIDs, apiKey.ActingUser.ID)
		}
	}

	// Loads acting users
	actingUsers, err := uc.userRepo.ListByIDs(ctx, uc.db, userIDs, bunex.SelectWithDeleted())
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	userMap := entityutil.SliceToIDMap(actingUsers)

	resp, err := apikeydto.TransformAPIKeys(settings, userMap)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.ListAPIKeyResp{
		Meta: &basedto.Meta{Page: paging},
		Data: resp,
	}, nil
}
