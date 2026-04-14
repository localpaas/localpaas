package configfileuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/configfileuc/configfiledto"
)

func (uc *UC) CreateConfigFile(
	ctx context.Context,
	auth *basedto.Auth,
	req *configfiledto.CreateConfigFileReq,
) (*configfiledto.CreateConfigFileResp, error) {
	req.Type = currentSettingType
	resp, err := uc.CreateSetting(ctx, &req.CreateSettingReq, &settings.CreateSettingData{
		VerifyingName: req.Name,
		Version:       currentSettingVersion,
		PrepareCreation: func(
			ctx context.Context,
			db database.Tx,
			data *settings.CreateSettingData,
			pData *settings.PersistingSettingCreationData,
		) error {
			configFile := req.ToEntity()
			if data.ScopeApp != nil {
				// Create a config in docker swarm
				err := uc.AppService.CreateSwarmConfig(ctx, db, data.ScopeApp, configFile)
				if err != nil {
					return apperrors.Wrap(err)
				}
			}

			err := pData.Setting.SetData(configFile)
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &configfiledto.CreateConfigFileResp{
		Data: resp.Data,
	}, nil
}
