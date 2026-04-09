package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentImageBuildSettingsVersion = 1

	defaultCPUPeriod = 100000
)

var _ = registerSettingParser(base.SettingTypeImageBuildSettings, &imageBuildSettingsParser{})

type imageBuildSettingsParser struct {
}

func (s *imageBuildSettingsParser) New() SettingData {
	return &ImageBuildSettings{}
}

type ImageBuildSettings struct {
	Resources *ImageBuildSettingResources `json:"resources"`
	NoCache   bool                        `json:"noCache,omitempty"`
	NoVerbose bool                        `json:"noVerbose,omitempty"`
}

type ImageBuildSettingResources struct {
	CPUs      int32 `json:"cpus"`
	MemMB     int64 `json:"memMB"`
	MemSwapMB int64 `json:"memSwapMB,omitempty"`
	ShmSizeMB int64 `json:"shmSizeMB,omitempty"`
}

// CPUsAsPeriodAndQuota calculates CPU period and quota from CPUs
// Ref: https://docs.docker.com/engine/containers/resource_constraints
func (s *ImageBuildSettingResources) CPUsAsPeriodAndQuota() (int64, int64) {
	if s.CPUs == 0 {
		return 0, 0
	}
	return defaultCPUPeriod, int64(defaultCPUPeriod * s.CPUs)
}

func (s *ImageBuildSettings) GetType() base.SettingType {
	return base.SettingTypeImageBuildSettings
}

func (s *ImageBuildSettings) GetRefObjectIDs() *RefObjectIDs {
	return &RefObjectIDs{}
}

func (s *ImageBuildSettings) Migrate(setting *Setting) (hasChange bool, err error) {
	if setting.Version == CurrentImageBuildSettingsVersion {
		return false, nil
	}
	if setting.Version > CurrentImageBuildSettingsVersion {
		return false, apperrors.New(apperrors.ErrDataVerNewerThanSystemVer)
	}

	// TODO: add migration if we make any change

	setting.Version = CurrentImageBuildSettingsVersion
	setting.UpdateVer++
	setting.MustSetData(s)
	return true, nil
}

func (s *Setting) AsImageBuildSettings() (*ImageBuildSettings, error) {
	return parseSettingAs[*ImageBuildSettings](s)
}

func (s *Setting) MustAsImageBuildSettings() *ImageBuildSettings {
	return gofn.Must(s.AsImageBuildSettings())
}
