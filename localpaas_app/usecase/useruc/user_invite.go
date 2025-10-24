package useruc

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/userservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/ulid"
)

const (
	dashboardUserSignUpPath = "/auth/sign-up"
)

func (uc *UserUC) InviteUser(
	ctx context.Context,
	auth *basedto.Auth,
	req *userdto.InviteUserReq,
) (*userdto.InviteUserResp, error) {
	inviteData := &userInviteData{}
	err := uc.loadUserInviteData(ctx, uc.db, req, inviteData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &userservice.PersistingUserData{}
	uc.preparePersistingUserInviteData(req, inviteData, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.userService.PersistUserData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// TODO: handle if req.SendInviteEmail = true

	return &userdto.InviteUserResp{
		Data: &userdto.InviteUserDataResp{
			InviteLink: inviteData.InviteLink,
		},
	}, nil
}

type userInviteData struct {
	User       *entity.User
	InviteLink string
}

func (uc *UserUC) loadUserInviteData(
	ctx context.Context,
	db database.IDB,
	req *userdto.InviteUserReq,
	data *userInviteData,
) error {
	user, err := uc.userRepo.GetByEmail(ctx, db, req.Email)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if user != nil && user.Status != base.UserStatusPending {
		return apperrors.NewAlreadyExist("User").
			WithMsgLog("user '%s' already exists", req.Email)
	}

	if user == nil {
		user = &entity.User{
			ID:        gofn.Must(ulid.NewStringULID()),
			Email:     req.Email,
			CreatedAt: time.Now(),
		}
	} else { //nolint
		// TODO: remove all old accesses
	}
	data.User = user

	// Generate invite token
	inviteToken, err := uc.userService.GenerateUserInviteToken(user.ID)
	if err != nil {
		return apperrors.Wrap(err)
	}
	signupLink := gofn.Must(url.JoinPath(config.Current.App.BaseURL, dashboardUserSignUpPath)) +
		fmt.Sprintf("?token=%s", inviteToken)
	data.InviteLink = signupLink

	return nil
}

func (uc *UserUC) preparePersistingUserInviteData(
	req *userdto.InviteUserReq,
	data *userInviteData,
	persistingData *userservice.PersistingUserData,
) {
	timeNow := timeutil.NowUTC()
	user := data.User

	user.Role = req.Role
	user.Status = base.UserStatusPending
	user.SecurityOption = req.SecurityOption
	user.AccessExpireAt = req.AccessExpiration
	user.UpdatedAt = timeNow

	persistingData.UpsertingUsers = append(persistingData.UpsertingUsers, user)

	uc.preparePersistingUserModuleAccesses(user, req.ModuleAccesses, timeNow, persistingData)
	uc.preparePersistingUserProjectAccesses(user, req.ProjectAccesses, timeNow, persistingData)
}

func (uc *UserUC) preparePersistingUserModuleAccesses(
	user *entity.User,
	moduleReqs basedto.ModuleAccessSliceReq,
	timeNow time.Time,
	persistingData *userservice.PersistingUserData,
) {
	for _, moduleReq := range moduleReqs {
		persistingData.UpsertingAccesses = append(persistingData.UpsertingAccesses,
			&entity.ACLPermission{
				SubjectType:  base.SubjectTypeUser,
				SubjectID:    user.ID,
				ResourceType: base.ResourceTypeModule,
				ResourceID:   moduleReq.ID,
				Actions:      moduleReq.Access,
				CreatedAt:    timeNow,
				UpdatedAt:    timeNow,
			})
	}
}

func (uc *UserUC) preparePersistingUserProjectAccesses(
	user *entity.User,
	projectReqs basedto.ObjectAccessSliceReq,
	timeNow time.Time,
	persistingData *userservice.PersistingUserData,
) {
	for _, projectReq := range projectReqs {
		persistingData.UpsertingAccesses = append(persistingData.UpsertingAccesses,
			&entity.ACLPermission{
				SubjectType:  base.SubjectTypeUser,
				SubjectID:    user.ID,
				ResourceType: base.ResourceTypeProject,
				ResourceID:   projectReq.ID,
				Actions:      projectReq.Access,
				CreatedAt:    timeNow,
				UpdatedAt:    timeNow,
			})
	}
}
