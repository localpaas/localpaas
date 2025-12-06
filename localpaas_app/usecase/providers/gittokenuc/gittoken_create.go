package gittokenuc

import (
	"context"
	"errors"

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
)

func (uc *GitTokenUC) CreateGitToken(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.CreateGitTokenReq,
) (*gittokendto.CreateGitTokenResp, error) {
	tokenData := &createGitTokenData{}
	err := uc.loadGitTokenData(ctx, uc.db, req, tokenData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingGitTokenData{}
	uc.preparePersistingGitToken(req, tokenData, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &gittokendto.CreateGitTokenResp{
		Data: &basedto.ObjectIDResp{ID: createdItem.ID},
	}, nil
}

type createGitTokenData struct {
}

func (uc *GitTokenUC) loadGitTokenData(
	ctx context.Context,
	db database.IDB,
	req *gittokendto.CreateGitTokenReq,
	_ *createGitTokenData,
) error {
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeGitToken, req.Name, false)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("GitToken").
			WithMsgLog("git token '%s' already exists", req.Name)
	}

	return nil
}

type persistingGitTokenData struct {
	settingservice.PersistingSettingData
}

func (uc *GitTokenUC) preparePersistingGitToken(
	req *gittokendto.CreateGitTokenReq,
	_ *createGitTokenData,
	persistingData *persistingGitTokenData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeGitToken,
		Status:    base.SettingStatusActive,
		Kind:      string(req.TokenType),
		Name:      req.Name,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		ExpireAt:  req.ExpireAt,
	}

	githubApp := &entity.GitToken{
		User:  req.User,
		Token: entity.NewEncryptedField(req.Token),
	}
	setting.MustSetData(githubApp)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}

func (uc *GitTokenUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingGitTokenData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
