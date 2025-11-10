package sessionuc

import (
	"context"
	"errors"
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

const (
	externalAvatarURLMaxLen = 2000
)

func (uc *SessionUC) CreateOAuthSession(
	ctx context.Context,
	req *sessiondto.CreateOAuthSessionReq,
) (*sessiondto.CreateOAuthSessionResp, error) {
	oauthUser := req.User
	email := gofn.Coalesce(oauthUser.Email, oauthUser.UserID)
	dbUser, err := uc.userRepo.GetByEmail(ctx, uc.db, email)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, apperrors.Wrap(err)
	}

	fullName := strings.Join(gofn.ToSliceSkippingZero(oauthUser.FirstName, oauthUser.LastName), " ")
	fullName = gofn.Coalesce(fullName, oauthUser.NickName, oauthUser.UserID)

	if dbUser == nil {
		// User makes the first login to our service
		timeNow := timeutil.NowUTC()
		dbUser = &entity.User{
			Email:          email,
			FullName:       fullName,
			Photo:          gofn.If(len(oauthUser.AvatarURL) < externalAvatarURLMaxLen, oauthUser.AvatarURL, ""), //nolint
			Role:           base.UserRoleMember,
			Status:         base.UserStatusActive,
			SecurityOption: base.UserSecurityEnforceSSO,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}
		if err = uc.userRepo.Insert(ctx, uc.db, dbUser); err != nil {
			return nil, apperrors.New(err).WithMsgLog("failed to insert new user to db: %s", dbUser.Email)
		}
	} else {
		// Make synchronization as info of user may be changed in the IdP system
		var updateCols []string
		if strings.HasPrefix(dbUser.Photo, "http") && len(oauthUser.AvatarURL) < externalAvatarURLMaxLen {
			dbUser.Photo = oauthUser.AvatarURL
			updateCols = append(updateCols, "photo")
		}
		// Saves the user (NOTE: ignore the error as it may not be important)
		_ = uc.userRepo.Update(ctx, uc.db, dbUser, bunex.UpdateColumns(updateCols...))
	}

	sessionResp, err := uc.createSession(ctx, &sessiondto.BaseCreateSessionReq{User: dbUser})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sessiondto.CreateOAuthSessionResp{
		BaseCreateSessionResp: *sessionResp,
	}, nil
}
