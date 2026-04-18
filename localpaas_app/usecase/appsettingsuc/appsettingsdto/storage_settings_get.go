package appsettingsdto

import (
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/osutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/unit"
)

type GetAppStorageSettingsReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`
}

func NewGetAppStorageSettingsReq() *GetAppStorageSettingsReq {
	return &GetAppStorageSettingsReq{}
}

func (req *GetAppStorageSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAppStorageSettingsResp struct {
	Meta *basedto.Meta        `json:"meta"`
	Data *StorageSettingsResp `json:"data"`
}

type StorageSettingsResp struct {
	Mounts []*Mount `json:"mounts"`

	UpdateVer int `json:"updateVer"`
}

type Mount struct {
	Type           mount.Type        `json:"type"`
	Source         string            `json:"source,omitempty"`
	Target         string            `json:"target"`
	ReadOnly       bool              `json:"readOnly,omitempty"`
	Consistency    mount.Consistency `json:"consistency,omitempty"`
	BindOptions    *BindOptions      `json:"bindOptions,omitempty"`
	VolumeOptions  *VolumeOptions    `json:"volumeOptions,omitempty"`
	TmpfsOptions   *TmpfsOptions     `json:"tmpfsOptions,omitempty"`
	ClusterOptions *ClusterOptions   `json:"clusterOptions,omitempty"`
}

type BindOptions struct {
	Propagation            mount.Propagation `json:"propagation"`
	NonRecursive           bool              `json:"nonRecursive"`
	CreateMountpoint       bool              `json:"createMountpoint"`
	ReadOnlyNonRecursive   bool              `json:"readOnlyNonRecursive"`
	ReadOnlyForceRecursive bool              `json:"readOnlyForceRecursive"`
}

type VolumeOptions struct {
	NoCopy       bool              `json:"noCopy"`
	Labels       map[string]string `json:"labels"`
	Subpath      string            `json:"subpath"`
	DriverConfig *VolumeDriver     `json:"driverConfig"`
}

type VolumeDriver struct {
	Name    string            `json:"name"`
	Options map[string]string `json:"options"`
}

type TmpfsOptions struct {
	Size    unit.DataSize   `json:"size" copy:"SizeBytes"`
	Mode    osutil.FileMode `json:"mode"`
	Options [][]string      `json:"options"`
}

type ClusterOptions struct {
	VolumeOptions
}

func TransformStorageSettings(
	service *swarm.Service,
) (resp *StorageSettingsResp, err error) {
	spec := &service.Spec
	resp = &StorageSettingsResp{
		UpdateVer: int(service.Version.Index), //nolint:gosec
	}

	resp.Mounts, err = TransformStorageMounts(spec.TaskTemplate.ContainerSpec.Mounts)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}

func TransformStorageMount(mount *mount.Mount) (resp *Mount, err error) {
	if err = copier.Copy(&resp, mount); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

func TransformStorageMounts(mounts []mount.Mount) ([]*Mount, error) {
	resp := make([]*Mount, 0, len(mounts))
	for _, mnt := range mounts {
		itemResp, err := TransformStorageMount(&mnt)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, itemResp)
	}
	return resp, nil
}
