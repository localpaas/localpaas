package gittokenuc

import (
	"context"
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
)

func (uc *GitTokenUC) UpdateGitToken(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.UpdateGitTokenReq,
) (*gittokendto.UpdateGitTokenResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		tokenData := &updateGitTokenData{}
		err := uc.loadGitTokenDataForUpdate(ctx, db, req, tokenData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingGitTokenData{}
		uc.prepareUpdatingGitToken(req.GitTokenBaseReq, tokenData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gittokendto.UpdateGitTokenResp{}, nil
}

type updateGitTokenData struct {
	Setting *entity.Setting
}

func (uc *GitTokenUC) loadGitTokenDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *gittokendto.UpdateGitTokenReq,
	data *updateGitTokenData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeGitToken, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting

	// If name changes, validate the new one
	if req.Name != "" && !strings.EqualFold(setting.Name, req.Name) {
		conflictSetting, _ := uc.settingRepo.GetByName(ctx, db, base.SettingTypeGitToken, req.Name, false)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist("GitToken").
				WithMsgLog("git token '%s' already exists", conflictSetting.Name)
		}
	}

	return nil
}

func (uc *GitTokenUC) prepareUpdatingGitToken(
	req *gittokendto.GitTokenBaseReq,
	data *updateGitTokenData,
	persistingData *persistingGitTokenData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	setting.UpdateVer++
	setting.UpdatedAt = timeNow
	setting.Kind = gofn.Coalesce(string(req.Kind), setting.Kind)
	setting.Name = gofn.Coalesce(req.Name, setting.Name)
	setting.ExpireAt = req.ExpireAt

	githubApp := &entity.GitToken{
		User:  req.User,
		Token: entity.NewEncryptedField(req.Token),
	}
	setting.MustSetData(githubApp)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
