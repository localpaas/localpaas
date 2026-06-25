package apppreviewuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/apppreviewuc/apppreviewdto"
)

func (uc *UC) PrepareCreatePreview(
	ctx context.Context,
	auth *basedto.Auth,
	req *apppreviewdto.PrepareCreatePreviewReq,
) (_ *apppreviewdto.PrepareCreatePreviewResp, err error) {
	app, err := uc.appService.LoadApp(ctx, uc.db, req.ProjectID, req.AppID, true, true,
		bunex.SelectExcludeColumns(entity.AppDefaultExcludeColumns...),
		bunex.SelectRelation("Project",
			bunex.SelectExcludeColumns(entity.ProjectDefaultExcludeColumns...),
		),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type = ?", base.SettingTypeAppDeployment),
		),
	)
	if err != nil {
		return nil, apperrors.New(err)
	}

	deploymentSetting := app.GetSettingByType(base.SettingTypeAppDeployment)
	if deploymentSetting == nil {
		return nil, apperrors.NewNotFound("Deployment settings")
	}
	deploymentSettings := deploymentSetting.MustAsAppDeploymentSettings()
	repoSource := deploymentSettings.RepoSource
	if deploymentSettings.ActiveMethod != base.DeploymentMethodRepo || repoSource == nil {
		return nil, apperrors.New(apperrors.ErrDeploymentMethodRepoRequired)
	}

	refObjects, err := uc.settingService.LoadReferenceObjects(ctx, uc.db, app.GetObjectScope(),
		true, true, app.Settings...)
	if err != nil {
		return nil, apperrors.New(err)
	}

	respData := &apppreviewdto.PrepareCreatePreviewDataResp{}
	if repoSource.Credentials.ID != "" {
		respData.RepoCredentials = &basedto.ObjectIDResp{ID: repoSource.Credentials.ID}
	}

	credSetting := refObjects.RefSettings[repoSource.Credentials.ID]
	if credSetting != nil {
		if credSetting.Type == base.SettingTypeGithubApp || credSetting.Type == base.SettingTypeAccessToken {
			respData.CanListBranches = true
			respData.CanListPullRequests = true
		}
	}

	return &apppreviewdto.PrepareCreatePreviewResp{
		Data: respData,
	}, nil
}
