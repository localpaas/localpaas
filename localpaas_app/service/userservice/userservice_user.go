package userservice

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/entityutil"
)

func (s *userService) LoadUser(
	ctx context.Context,
	db database.IDB,
	userID string,
) (*entity.User, error) {
	userMap, err := s.LoadUsers(ctx, db, []string{userID})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(userMap) == 0 {
		return nil, apperrors.NewNotFound("User")
	}
	return userMap[userID], nil
}

func (s *userService) LoadUsers(
	ctx context.Context,
	db database.IDB,
	userIDs []string,
) (userMap map[string]*entity.User, err error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	userIDs = gofn.ToSet(userIDs)
	users, err := s.userRepo.ListByIDs(ctx, db, userIDs)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	userMap = entityutil.SliceToIDMap(users)
	for _, userID := range userIDs {
		user := userMap[userID]
		if user == nil {
			return nil, apperrors.NewNotFound("User").
				WithMsgLog("user %s not found", userID)
		}
		if user.Status != base.UserStatusActive {
			return nil, apperrors.New(apperrors.ErrUserUnavailable).
				WithMsgLog("user %s is not active", userID)
		}
		if user.IsAccessExpired() {
			return nil, apperrors.New(apperrors.ErrUserUnavailable).
				WithMsgLog("user access expired at: %v", user.AccessExpireAt)
		}
		if user.SecurityOption == base.UserSecurityPassword2FA && user.TotpSecret == "" {
			return nil, apperrors.New(apperrors.ErrUserNotCompleteMFASetup).
				WithMsgLog("user %s hasn't completed the MFA setup", userID)
		}
	}

	return userMap, nil
}

func (s *userService) LoadUserByEmail(
	ctx context.Context,
	db database.IDB,
	email string,
) (*entity.User, error) {
	userMap, err := s.LoadUsersByEmails(ctx, db, []string{email})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(userMap) == 0 {
		return nil, apperrors.NewNotFound("User")
	}
	return userMap[email], nil
}

func (s *userService) LoadUsersByEmails(
	ctx context.Context,
	db database.IDB,
	emails []string,
) (userMap map[string]*entity.User, err error) {
	if len(emails) == 0 {
		return nil, nil
	}
	users, err := s.userRepo.ListByEmails(ctx, db, emails)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	userMap = make(map[string]*entity.User, len(users))
	for _, user := range users {
		userMap[user.Email] = user
	}

	for _, email := range emails {
		user := userMap[email]
		if user == nil {
			return nil, apperrors.NewNotFound("User").
				WithMsgLog("user '%s' not found", email)
		}
		if user.Status != base.UserStatusActive {
			return nil, apperrors.New(apperrors.ErrUserUnavailable).
				WithMsgLog("user '%s' is not active", email)
		}
		if user.IsAccessExpired() {
			return nil, apperrors.New(apperrors.ErrUserUnavailable).
				WithMsgLog("user access expired at: %v", user.AccessExpireAt)
		}
		if user.SecurityOption == base.UserSecurityPassword2FA && user.TotpSecret == "" {
			return nil, apperrors.New(apperrors.ErrUserNotCompleteMFASetup).
				WithMsgLog("user '%s' hasn't completed the MFA setup", email)
		}
	}

	return userMap, nil
}
