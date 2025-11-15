package base

type VolumeDriver string

const (
	VolumeDriverLocal     = VolumeDriver("local")
	VolumeDriverSeaweedFs = VolumeDriver("seaweedfs")
)

var (
	AllVolumeDrivers = []VolumeDriver{VolumeDriverLocal, VolumeDriverSeaweedFs}
)

type VolumeType string

const (
	VolumeTypeVolume = VolumeType("volume")
	VolumeTypeNfs    = VolumeType("nfs")
)

var (
	AllVolumeTypes = []VolumeType{VolumeTypeVolume, VolumeTypeNfs}
)

type VolumeScope string

const (
	VolumeScopeGlobal = VolumeScope("global")
	VolumeScopeLocal  = VolumeScope("local")
)

var (
	AllVolumeScopes = []VolumeScope{VolumeScopeGlobal, VolumeScopeLocal}
)
