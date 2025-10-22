package appuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) GetAppEnvVars(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.GetAppEnvVarsReq,
) (*appdto.GetAppEnvVarsResp, error) {
	app, err := uc.appRepo.GetByID(ctx, uc.db, req.ProjectID, req.AppID,
		bunex.SelectRelation("EnvVarsSettings"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	envVars, err := app.GetEnvVars()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := appdto.TransformAppEnvVars(envVars)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.GetAppEnvVarsResp{
		Data: resp,
	}, nil
}
