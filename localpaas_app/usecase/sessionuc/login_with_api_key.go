package sessionuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
)

func (uc *SessionUC) LoginWithAPIKey(
	ctx context.Context,
	req *sessiondto.LoginWithAPIKeyReq,
) (resp *sessiondto.LoginWithAPIKeyResp, err error) {
	apiKeySetting, err := uc.settingRepo.GetByName(ctx, uc.db, base.SettingTypeAPIKey, req.KeyID, false)
	if err != nil {
		return nil, uc.wrapSensitiveError(err)
	}
	if !apiKeySetting.IsActive() {
		return nil, uc.wrapSensitiveError(apperrors.ErrAPIKeyInvalid)
	}

	apiKey, err := apiKeySetting.ParseAPIKey()
	if err != nil {
		return nil, uc.wrapSensitiveError(err)
	}
	if apiKey == nil {
		return nil, uc.wrapSensitiveError(apperrors.ErrAPIKeyMismatched)
	}
	if err = apiKey.VerifyHash(req.SecretKey); err != nil {
		return nil, uc.wrapSensitiveError(err)
	}
	actingUserID := apiKeySetting.ObjectID

	dbUser, err := uc.userService.LoadUser(ctx, uc.db, actingUserID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Create a new session as login succeeds
	sessionData, err := uc.createSession(ctx, &sessiondto.BaseCreateSessionReq{
		User:         dbUser,
		IsAPIKey:     true,
		AccessAction: apiKey.AccessAction,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sessiondto.LoginWithAPIKeyResp{
		Data: &sessiondto.LoginWithAPIKeyDataResp{
			AccessToken:     sessionData.AccessToken,
			AccessTokenExp:  sessionData.AccessTokenExp,
			RefreshToken:    sessionData.RefreshToken,
			RefreshTokenExp: sessionData.RefreshTokenExp,
		},
	}, nil
}
