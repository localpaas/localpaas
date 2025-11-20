package docker

import (
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

const (
	UnitCPUNano = 1000 * 1000 * 1000
	UnitMemMB   = 1024 * 1024
	memMBMin    = 10
)

type ServiceSpec struct {
	Name                 string               `json:"name"`
	Labels               map[string]string    `json:"labels,omitempty"`
	Image                string               `json:"image"`
	Command              []string             `json:"command,omitempty"`
	ServiceMode          ServiceMode          `json:"serviceMode"`
	Replicas             uint64               `json:"replicas"`
	Hostname             string               `json:"hostname,omitempty"`
	Ports                []*PortConfig        `json:"ports,omitempty"`
	VolumeMounts         []*VolumeMount       `json:"volumeMounts,omitempty"`
	BindMounts           []*BindMount         `json:"bindMounts,omitempty"`
	Networks             []*NetworkAttachment `json:"networks,omitempty"`
	PlacementConstraints []string             `json:"placementConstraints,omitempty"`
	Healthcheck          *Healthcheck         `json:"healthcheck,omitempty"`
	ResourceReserved     *Resource            `json:"resourceReserved,omitempty"`
	ResourceLimit        *Resource            `json:"resourceLimit,omitempty"`
	Sysctls              map[string]string    `json:"sysctls,omitempty"`
	CapabilityAdd        []string             `json:"capabilityAdd,omitempty"`
	CapabilityDrop       []string             `json:"capabilityDrop,omitempty"`

	Env []string `json:"-"`

	// RestartPolicy *RestartPolicy
	// ForceUpdate is a counter that triggers an update even if no relevant
	// parameters have been changed.
	// ForceUpdate uint64
	// Runtime swarm.RuntimeType
}

type ServiceMode string

const (
	ServiceModeReplicated ServiceMode = "replicated"
	ServiceModeGlobal     ServiceMode = "global"
)

type NetworkAttachment struct {
	Target string `json:"target,omitempty"`
	// Aliases    []string
	// DriverOpts map[string]string
}

type PortConfig struct {
	Target    uint32 `json:"target,omitempty"`    // port inside the container
	Published uint32 `json:"published,omitempty"` // port on the swarm hosts
	Mode      string `json:"mode,omitempty"`      // mode in which port is published (ingress, host)
}

type VolumeMount struct {
	Source   string `json:"source,omitempty"`
	Target   string `json:"target,omitempty"`
	ReadOnly bool   `json:"readOnly,omitempty"`
}

type BindMount struct {
	Source   string `json:"source,omitempty"`
	Target   string `json:"target,omitempty"`
	ReadOnly bool   `json:"readOnly,omitempty"`
}

type Resource struct {
	CPUs     float64 `json:"cpus,omitempty"`
	MemoryMB int64   `json:"memoryMB,omitempty"`
}

type Healthcheck struct {
	Enabled bool `json:"enabled,omitempty"`

	// Test is the test to perform to check that the container is healthy.
	// An empty slice means to inherit the default.
	// The options are:
	// {} : inherit healthcheck
	// {"NONE"} : disable healthcheck
	// {"CMD", args...} : exec arguments directly
	// {"CMD-SHELL", command} : run command with system's default shell
	Mode    HealthcheckMode `json:"mode,omitempty"`
	Command string          `json:"command,omitempty"`

	// Zero means to inherit. Durations are expressed as integer nanoseconds.
	Interval      time.Duration `json:"interval,omitempty"`
	Timeout       time.Duration `json:"timeout,omitempty"`
	StartPeriod   time.Duration `json:"startPeriod,omitempty"`
	StartInterval time.Duration `json:"startInterval,omitempty"`

	// Retries is the number of consecutive failures needed to consider a container as unhealthy.
	// Zero means inherit.
	Retries int `json:"retries,omitempty"`
}

type HealthcheckMode string

const (
	HealthcheckModeInherit  = HealthcheckMode("")
	HealthcheckModeNone     = HealthcheckMode("NONE")
	HealthcheckModeCmd      = HealthcheckMode("CMD")
	HealthcheckModeCmdShell = HealthcheckMode("CMD_SHELL")
)

func (s *ServiceSpec) ToSwarmServiceSpec() (*swarm.ServiceSpec, error) {
	// Service mode
	var serviceMode swarm.ServiceMode
	switch s.ServiceMode {
	case ServiceModeReplicated:
		serviceMode = swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &s.Replicas,
			},
		}
	case ServiceModeGlobal:
		serviceMode = swarm.ServiceMode{
			Global: &swarm.GlobalService{},
		}
	default:
		return nil, tracerr.Wrap(ErrServiceModeNotSupported)
	}

	// Volumes
	var volumeMounts []mount.Mount
	for _, volumeMount := range s.VolumeMounts {
		volumeMounts = append(volumeMounts, mount.Mount{
			Type:     mount.TypeVolume,
			Source:   volumeMount.Source,
			Target:   volumeMount.Target,
			ReadOnly: volumeMount.ReadOnly,
		})
	}
	for _, volumeBind := range s.BindMounts {
		volumeMounts = append(volumeMounts, mount.Mount{
			Type:     mount.TypeBind,
			Source:   volumeBind.Source,
			Target:   volumeBind.Target,
			ReadOnly: volumeBind.ReadOnly,
		})
	}

	// Resources
	var resourceReserved *swarm.Resources
	var resourceLimit *swarm.Limit
	if s.ResourceReserved != nil {
		resourceReserved = &swarm.Resources{
			NanoCPUs:    int64(s.ResourceReserved.CPUs * UnitCPUNano),
			MemoryBytes: max(s.ResourceReserved.MemoryMB, memMBMin) * UnitMemMB,
		}
	}
	if s.ResourceLimit != nil {
		resourceLimit = &swarm.Limit{
			NanoCPUs:    int64(s.ResourceLimit.CPUs * UnitCPUNano),
			MemoryBytes: max(s.ResourceLimit.MemoryMB, memMBMin) * UnitMemMB,
		}
	}

	// healthcheck
	var healthcheck *container.HealthConfig
	if s.Healthcheck != nil && s.Healthcheck.Enabled {
		healthcheck = &container.HealthConfig{
			Test:          []string{string(s.Healthcheck.Mode), s.Healthcheck.Command},
			Interval:      s.Healthcheck.Interval,
			Timeout:       s.Healthcheck.Timeout,
			StartPeriod:   s.Healthcheck.StartPeriod,
			StartInterval: s.Healthcheck.StartInterval,
			Retries:       s.Healthcheck.Retries,
		}
	}

	return &swarm.ServiceSpec{
		// Annotations
		Annotations: swarm.Annotations{
			Name:   s.Name,
			Labels: s.Labels,
		},
		// Task spec
		TaskTemplate: swarm.TaskSpec{
			// Container spec
			ContainerSpec: &swarm.ContainerSpec{
				Image:    s.Image,
				Command:  s.Command,
				Hostname: s.Hostname,
				Env:      s.Env,
				Mounts:   volumeMounts,
				Privileges: &swarm.Privileges{
					NoNewPrivileges: true,
					AppArmor: &swarm.AppArmorOpts{
						Mode: swarm.AppArmorModeDefault,
					},
					Seccomp: &swarm.SeccompOpts{
						Mode: swarm.SeccompModeDefault,
					},
				},
				Healthcheck:    healthcheck,
				Sysctls:        s.Sysctls,
				CapabilityAdd:  s.CapabilityAdd,
				CapabilityDrop: s.CapabilityDrop,
			},
			// Networks
			Networks: gofn.MapSlice(s.Networks, func(net *NetworkAttachment) swarm.NetworkAttachmentConfig {
				return swarm.NetworkAttachmentConfig{
					Target: net.Target,
				}
			}),
			// Placement
			Placement: &swarm.Placement{
				Constraints: s.PlacementConstraints,
			},
			// Resources
			Resources: &swarm.ResourceRequirements{
				Reservations: resourceReserved,
				Limits:       resourceLimit,
			},
		},
		// Mode
		Mode: serviceMode,
		// Endpoint
		EndpointSpec: &swarm.EndpointSpec{
			Mode: swarm.ResolutionModeDNSRR,
			Ports: gofn.MapSlice(s.Ports, func(p *PortConfig) swarm.PortConfig {
				return swarm.PortConfig{
					TargetPort:    p.Target,
					PublishedPort: p.Published,
					PublishMode: gofn.If(p.Mode != "", swarm.PortConfigPublishMode(p.Mode), //nolint
						swarm.PortConfigPublishModeHost),
				}
			}),
		},
	}, nil
}
