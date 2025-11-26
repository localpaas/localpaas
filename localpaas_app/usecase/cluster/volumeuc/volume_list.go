package volumeuc

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/volumeuc/volumedto"
)

const (
	localpaasVolumeLabel = "localpaas.volume.managed"
)

func (uc *VolumeUC) ListVolume(
	ctx context.Context,
	auth *basedto.Auth,
	req *volumedto.ListVolumeReq,
) (*volumedto.ListVolumeResp, error) {
	volumes, err := uc.dockerManager.VolumeList(ctx, func(opts *volume.ListOptions) {
		if opts.Filters.Len() == 0 {
			opts.Filters = filters.NewArgs()
		}
		if !req.ListAll {
			opts.Filters.Add("label", localpaasVolumeLabel)
		}
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	filterVolumes := volumes.Volumes
	if req.Search != "" {
		keyword := strings.ToLower(req.Search)
		filterVolumes = gofn.Filter(filterVolumes, func(vol *volume.Volume) bool {
			return strings.Contains(vol.Name, keyword) || strings.Contains(vol.Mountpoint, keyword)
		})
	}
	if len(auth.AllowObjectIDs) > 0 {
		filterVolumes = gofn.Filter(filterVolumes, func(vol *volume.Volume) bool {
			volID := vol.Name
			if vol.ClusterVolume != nil {
				volID = vol.ClusterVolume.ID
			}
			return gofn.Contain(auth.AllowObjectIDs, volID)
		})
	}

	return &volumedto.ListVolumeResp{
		Meta: &basedto.Meta{Page: &basedto.PagingMeta{
			Offset: 0,
			Limit:  req.Paging.Limit,
			Total:  len(volumes.Volumes),
		}},
		Data: volumedto.TransformVolumes(filterVolumes, false),
	}, nil
}
