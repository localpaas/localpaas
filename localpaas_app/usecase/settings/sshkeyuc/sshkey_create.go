package sshkeyuc

import (
	"context"
	"errors"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc/sshkeydto"
)

func (uc *SSHKeyUC) CreateSSHKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *sshkeydto.CreateSSHKeyReq,
) (*sshkeydto.CreateSSHKeyResp, error) {
	sshKeyData := &createSSHKeyData{}
	err := uc.loadSSHKeyData(ctx, uc.db, req, sshKeyData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingSSHKeyData{}
	uc.preparePersistingSSHKey(req.SSHKeyBaseReq, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &sshkeydto.CreateSSHKeyResp{
		Data: &basedto.ObjectIDResp{ID: createdItem.ID},
	}, nil
}

type createSSHKeyData struct {
}

func (uc *SSHKeyUC) loadSSHKeyData(
	ctx context.Context,
	db database.IDB,
	req *sshkeydto.CreateSSHKeyReq,
	_ *createSSHKeyData,
) error {
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeSSHKey, req.Name, false)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("SSHKey").
			WithMsgLog("SSH key '%s' setting already exists", req.Name)
	}

	return nil
}

type persistingSSHKeyData struct {
	settingservice.PersistingSettingData
}

func (uc *SSHKeyUC) preparePersistingSSHKey(
	req *sshkeydto.SSHKeyBaseReq,
	persistingData *persistingSSHKeyData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeSSHKey,
		Status:    base.SettingStatusActive,
		Name:      req.Name,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	sshKey := &entity.SSHKey{
		PrivateKey: req.PrivateKey,
		Passphrase: req.Passphrase,
	}
	setting.MustSetData(sshKey.MustEncrypt())

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)

	uc.preparePersistingSSHKeyProjects(setting, req.ProjectAccesses, timeNow, persistingData)
}

func (uc *SSHKeyUC) preparePersistingSSHKeyProjects(
	setting *entity.Setting,
	projectReqs []*sshkeydto.SSHKeyProjectAccessReq,
	timeNow time.Time,
	persistingData *persistingSSHKeyData,
) {
	for _, projectReq := range projectReqs {
		persistingData.UpsertingAccesses = append(persistingData.UpsertingAccesses,
			&entity.ACLPermission{
				SubjectType:  base.SubjectTypeProject,
				SubjectID:    projectReq.ID,
				ResourceType: base.ResourceTypeSSHKey,
				ResourceID:   setting.ID,
				Actions:      entity.AccessActions{Read: projectReq.Allowed},
				CreatedAt:    timeNow,
				UpdatedAt:    timeNow,
			})
		uc.preparePersistingSSHKeyApps(setting, projectReq.AppAccesses, timeNow, persistingData)
	}
}

func (uc *SSHKeyUC) preparePersistingSSHKeyApps(
	setting *entity.Setting,
	appReqs []*sshkeydto.SSHKeyAppAccessReq,
	timeNow time.Time,
	persistingData *persistingSSHKeyData,
) {
	for _, appReq := range appReqs {
		persistingData.UpsertingAccesses = append(persistingData.UpsertingAccesses,
			&entity.ACLPermission{
				SubjectType:  base.SubjectTypeApp,
				SubjectID:    appReq.ID,
				ResourceType: base.ResourceTypeSSHKey,
				ResourceID:   setting.ID,
				Actions:      entity.AccessActions{Read: appReq.Allowed},
				CreatedAt:    timeNow,
				UpdatedAt:    timeNow,
			})
	}
}

func (uc *SSHKeyUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingSSHKeyData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
