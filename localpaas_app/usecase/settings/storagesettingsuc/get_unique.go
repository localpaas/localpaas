package storagesettingsuc

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types/volume"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/storagesettingsuc/storagesettingsdto"
)

func (uc *UC) GetUniqueStorageSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *storagesettingsdto.GetUniqueStorageSettingsReq,
) (*storagesettingsdto.GetUniqueStorageSettingsResp, error) {
	req.Type = currentSettingType
	resp, err := uc.GetUniqueSetting(ctx, auth, &req.GetUniqueSettingReq, &settings.GetUniqueSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	storageSetting := resp.Data.MustAsStorageSettings()
	// Load ref cluster volumes as their IDs are different from their names
	var volumes []*volume.Volume
	//nolint:nestif
	if storageSetting.ClusterVolumeSettings != nil && len(storageSetting.ClusterVolumeSettings.Volumes) > 0 {
		if len(storageSetting.ClusterVolumeSettings.Volumes) == 1 {
			vol, _, err := uc.dockerManager.VolumeInspect(ctx, storageSetting.ClusterVolumeSettings.Volumes[0].ID)
			if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
				return nil, apperrors.Wrap(err)
			}
			if vol != nil {
				volumes = append(volumes, vol)
			}
		} else {
			volResp, err := uc.dockerManager.VolumeList(ctx)
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
			volumes = append(volumes, volResp.Volumes...)
		}
	}

	respData, err := storagesettingsdto.TransformStorageSettings(resp.Data, volumes)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &storagesettingsdto.GetUniqueStorageSettingsResp{
		Data: respData,
	}, nil
}
