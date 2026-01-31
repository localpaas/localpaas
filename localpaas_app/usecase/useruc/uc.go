package useruc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/emailservice"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
)

type UserUC struct {
	db           *database.DB
	userRepo     repository.UserRepo
	userService  userservice.UserService
	emailService emailservice.EmailService
}

func NewUserUC(
	db *database.DB,
	userRepo repository.UserRepo,
	userService userservice.UserService,
	emailService emailservice.EmailService,
) *UserUC {
	return &UserUC{
		db:           db,
		userRepo:     userRepo,
		userService:  userService,
		emailService: emailService,
	}
}
