package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentImageBuildVersion = 1
)

var _ = registerSettingParser(base.SettingTypeImageBuild, &imageBuildParser{})

type imageBuildParser struct {
}

func (s *imageBuildParser) New() SettingData {
	return &ImageBuild{}
}

type ImageBuild struct {
	Resources *ImageBuildResources `json:"resources"`
}

type ImageBuildResources struct {
	CPUs  uint32 `json:"cpus"`
	MemMB uint64 `json:"memMB"`
}

func (s *ImageBuild) GetType() base.SettingType {
	return base.SettingTypeImageBuild
}

func (s *ImageBuild) GetRefObjectIDs() *RefObjectIDs {
	return &RefObjectIDs{}
}

func (s *Setting) AsImageBuild() (*ImageBuild, error) {
	return parseSettingAs[*ImageBuild](s)
}

func (s *Setting) MustAsImageBuild() *ImageBuild {
	return gofn.Must(s.AsImageBuild())
}
