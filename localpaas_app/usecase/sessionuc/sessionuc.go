package sessionuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/redisrepository"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type SessionUC struct {
	db                     *database.DB
	userRepo               repository.UserRepo
	loginTrustedDeviceRepo repository.LoginTrustedDeviceRepo
	userTokenRepo          redisrepository.UserTokenRepo
	mfaPasscodeRepo        redisrepository.MFAPasscodeRepo
	userService            userservice.UserService
	permissionManager      permission.Manager
}

func NewSessionUC(
	db *database.DB,
	userRepo repository.UserRepo,
	loginTrustedDeviceRepo repository.LoginTrustedDeviceRepo,
	userTokenRepo redisrepository.UserTokenRepo,
	mfaPasscodeRepo redisrepository.MFAPasscodeRepo,
	userService userservice.UserService,
	permissionManager permission.Manager,
) *SessionUC {
	return &SessionUC{
		db:                     db,
		userRepo:               userRepo,
		loginTrustedDeviceRepo: loginTrustedDeviceRepo,
		userTokenRepo:          userTokenRepo,
		mfaPasscodeRepo:        mfaPasscodeRepo,
		userService:            userService,
		permissionManager:      permissionManager,
	}
}
