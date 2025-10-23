package sessionuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
	"github.com/localpaas/localpaas/pkg/randtoken"
)

const (
	hashingKeyLen    = 32
	hashingIteration = 1
)

func (uc *SessionUC) LoginWithAPIKey(
	ctx context.Context,
	req *sessiondto.LoginWithAPIKeyReq,
) (resp *sessiondto.LoginWithAPIKeyResp, err error) {
	apiKeySetting, err := uc.settingRepo.GetByName(ctx, uc.db, base.SettingTypeAPIKey, req.KeyID)
	if err != nil {
		return nil, uc.wrapSensitiveError(err)
	}
	apiKey, err := apiKeySetting.ParseAPIKey()
	if err != nil {
		return nil, uc.wrapSensitiveError(err)
	}
	if apiKey == nil || apiKey.ActingUser.ID == "" {
		return nil, uc.wrapSensitiveError(apperrors.ErrAPIKeyMismatched)
	}
	if !randtoken.VerifyHashHex(req.SecretKey, apiKey.SecretKey, apiKey.Salt, hashingKeyLen, hashingIteration) {
		return nil, uc.wrapSensitiveError(apperrors.ErrAPIKeyMismatched)
	}

	dbUser, err := uc.userService.LoadUser(ctx, uc.db, apiKey.ActingUser.ID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Create a new session as login succeeds
	sessionData, err := uc.createSession(ctx, &sessiondto.BaseCreateSessionReq{User: dbUser})
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
