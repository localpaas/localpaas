package userservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type UserService interface {
	LoadUser(ctx context.Context, db database.IDB, userID string) (*entity.User, error)
	LoadUsers(ctx context.Context, db database.IDB, userIDs []string) (
		userMap map[string]*entity.User, err error)
	LoadUserByEmail(ctx context.Context, db database.IDB, email string) (*entity.User, error)
	LoadUsersByEmails(ctx context.Context, db database.IDB, emails []string) (
		userMap map[string]*entity.User, err error)

	ChangePassword(user *entity.User, newPassword, currPassword string) error
	VerifyPassword(user *entity.User, password string) error
	CheckPasswordStrength(password string) error

	GenerateMFAToken(userID string, mfaType base.MFAType, trustedDeviceID string) (string, error)
}

func NewUserService(
	userRepo repository.UserRepo,
) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

type userService struct {
	userRepo repository.UserRepo
}
