package sessionuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
)

func (uc *SessionUC) GetLoginOptions(
	ctx context.Context,
	req *sessiondto.GetLoginOptionsReq,
) (*sessiondto.GetLoginOptionsResp, error) {
	settings, _, err := uc.settingRepo.List(ctx, uc.db, nil,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeOAuth),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectColumns("id", "kind", "name"),
		bunex.SelectOrder("kind"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	var resp []*sessiondto.LoginOptionResp
	for _, setting := range settings {
		oauthType := base.OAuthType(setting.Kind)
		resp = append(resp, &sessiondto.LoginOptionResp{
			Type:    oauthType,
			Name:    setting.Name,
			AuthURL: "/_/auth/sso/" + string(oauthType),
		})
	}

	return &sessiondto.GetLoginOptionsResp{
		Data: resp,
	}, nil
}
