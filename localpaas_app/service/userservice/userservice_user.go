package userservice

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/entityutil"
)

func (s *userService) LoadUsers(ctx context.Context, db database.IDB, userIDs []string) (
	userMap map[string]*entity.User, err error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	userIDs = gofn.ToSet(userIDs)
	users, err := s.userRepo.ListByIDs(ctx, db, userIDs)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	userMap = entityutil.SliceToIDMap(users)

	return userMap, nil
}

func (s *userService) LoadUsersByEmails(ctx context.Context, db database.IDB,
	emails []string) (userMap map[string]*entity.User, err error) {
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

	return userMap, nil
}
