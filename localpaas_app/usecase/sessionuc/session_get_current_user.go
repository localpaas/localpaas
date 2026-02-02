package sessionuc

import (
	"context"
	"errors"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/jwtsession"
)

func (uc *SessionUC) GetCurrentUser(ctx context.Context, jwt string) (*basedto.User, error) {
	authClaims := &jwtsession.AuthClaims{}
	err := jwtsession.ParseToken(jwt, authClaims)
	if err != nil {
		if errors.Is(err, jwtsession.ErrTokenExpired) {
			return nil, apperrors.New(apperrors.ErrSessionJWTExpired).WithCause(err)
		}
		return nil, apperrors.New(apperrors.ErrSessionJWTInvalid).WithCause(err)
	}

	// Make sure the token is marked `existing` in redis
	if err = uc.userTokenRepo.Exist(ctx, authClaims.UserID, authClaims.UID); err != nil {
		return nil, apperrors.New(apperrors.ErrSessionJWTInvalid).WithCause(err)
	}

	user, err := uc.userRepo.GetByID(ctx, uc.db, authClaims.UserID,
		bunex.SelectExcludeColumns(entity.UserDefaultExcludeColumns...),
	)
	if err != nil {
		return nil, apperrors.New(err)
	}

	return &basedto.User{User: user, AuthClaims: authClaims}, nil
}
