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
		bunex.SelectColumns("id", "name"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp := &sessiondto.LoginOptionsResp{
		AllowLoginWithGitHub: false,
		AllowLoginWithGitLab: false,
	}

	for _, setting := range settings {
		switch setting.Name {
		case "github":
			resp.AllowLoginWithGitHub = true
		case "gitlab":
			resp.AllowLoginWithGitLab = true
		case "google":
			resp.AllowLoginWithGoogle = true
		}
	}

	return &sessiondto.GetLoginOptionsResp{
		Data: resp,
	}, nil
}
