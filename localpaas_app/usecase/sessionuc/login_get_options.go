package sessionuc

import (
	"context"

	"github.com/tiendc/gofn"

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
		bunex.SelectColumns("id", "name"),
		bunex.SelectOrder("name"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	var resp []*sessiondto.LoginOptionResp
	mapDisplayName := map[base.OAuthType]string{
		base.OAuthTypeGitlabCustom: "Our Gitlab",
	}

	for _, setting := range settings {
		oauthType := base.OAuthType(setting.Name)
		displayName := mapDisplayName[oauthType]
		if displayName == "" {
			displayName = gofn.StringToUpper1stLetter(setting.Name)
		}
		resp = append(resp, &sessiondto.LoginOptionResp{
			Type:    oauthType,
			Name:    displayName,
			AuthURL: "/_/auth/sso/" + string(oauthType),
		})
	}

	return &sessiondto.GetLoginOptionsResp{
		Data: resp,
	}, nil
}
