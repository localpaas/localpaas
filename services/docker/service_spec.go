package docker

import (
	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

const (
	UnitCPUNano = 1000 * 1000 * 1000
	UnitMemMB   = 1024 * 1024
)

type ServiceModeSpec struct {
	Mode                ServiceMode `json:"mode,omitempty"`
	ServiceReplicas     *uint64     `json:"serviceReplicas,omitempty"`
	JobMaxConcurrent    *uint64     `json:"jobMaxConcurrent,omitempty"`
	JobTotalCompletions *uint64     `json:"jobTotalCompletions,omitempty"`
}

type ServiceMode string

const (
	ServiceModeReplicated    ServiceMode = "replicated"
	ServiceModeReplicatedJob ServiceMode = "replicated-job"
	ServiceModeGlobal        ServiceMode = "global"
	ServiceModeGlobalJob     ServiceMode = "global-job"
)

type TaskSpec struct {
	Networks      []*NetworkAttachment  `json:"networks,omitempty"`
	Resources     *ResourceRequirements `json:"resources,omitempty"`
	Placement     *Placement            `json:"placement,omitempty"`
	RestartPolicy *RestartPolicy        `json:"restartPolicy,omitempty"`
}

type ContainerSpec struct {
	Labels           map[string]string  `json:"labels,omitempty"`
	Image            *string            `json:"image,omitempty"`
	Command          *string            `json:"command,omitempty"`
	WorkingDir       *string            `json:"workingDir,omitempty"`
	Hostname         *string            `json:"hostname,omitempty"`
	User             *string            `json:"user,omitempty"`
	Groups           []string           `json:"groups,omitempty"`
	StopSignal       *string            `json:"stopSignal,omitempty"`
	TTY              *bool              `json:"tty,omitempty"`
	OpenStdin        *bool              `json:"openStdin,omitempty"`
	ReadOnly         *bool              `json:"readOnly,omitempty"`
	StopGracePeriod  *timeutil.Duration `json:"stopGracePeriod,omitempty"` // e.g. 5s, 1m
	HostsFileEntries []*HostsFileEntry  `json:"hostsFileEntries,omitempty"`
	Ulimits          []*Ulimit          `json:"ulimits,omitempty"`
	Sysctls          map[string]string  `json:"sysctls,omitempty"`
	CapabilityAdd    []string           `json:"capabilityAdd,omitempty"`
	CapabilityDrop   []string           `json:"capabilityDrop,omitempty"`
	EnableGPU        *bool              `json:"enableGPU,omitempty"`
	Healthcheck      *Healthcheck       `json:"healthcheck,omitempty"`
	VolumeMounts     []*VolumeMount     `json:"volumeMounts,omitempty"`
	BindMounts       []*BindMount       `json:"bindMounts,omitempty"`
}

type NetworkAttachment struct {
	Target  string   `json:"target,omitempty"`
	Aliases []string `json:"aliases,omitempty"`
	// DriverOpts map[string]string `json:"driverOpts,omitempty"`
}

type EndpointSpec struct {
	Mode  swarm.ResolutionMode `json:"mode,omitempty"`
	Ports []*PortConfig        `json:"ports,omitempty"`
}

type PortConfig struct {
	Target      uint32                      `json:"target,omitempty"`    // port inside the container
	Published   uint32                      `json:"published,omitempty"` // port on the swarm hosts
	Protocol    swarm.PortConfigProtocol    `json:"protocol,omitempty"`
	PublishMode swarm.PortConfigPublishMode `json:"publishMode,omitempty"`
}

type HostsFileEntry struct {
	Address   string   `json:"address,omitempty"`
	Hostnames []string `json:"hostnames,omitempty"`
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

type Ulimit struct {
	Name string `json:"name,omitempty"`
	Hard int64  `json:"hard,omitempty"`
	Soft int64  `json:"soft,omitempty"`
}

type ResourceRequirements struct {
	Limits       *ResourceLimit `json:"limits,omitempty"`
	Reservations *Resources     `json:"reservations,omitempty"`
}

type Resources struct {
	CPUs             float64  `json:"cpus,omitempty"`
	MemoryMB         int64    `json:"memoryMB,omitempty"`
	GenericResources []string `json:"genericResources,omitempty"`
}

type ResourceLimit struct {
	CPUs     float64 `json:"cpus,omitempty"`
	MemoryMB int64   `json:"memoryMB,omitempty"`
	Pids     int64   `json:"pids,omitempty"`
}

type Placement struct {
	Constraints []string `json:"constraints,omitempty"`
	Preferences []string `json:"preferences,omitempty"`
}

type RestartPolicy struct {
	Condition   swarm.RestartPolicyCondition `json:"condition,omitempty"`
	Delay       *timeutil.Duration           `json:"delay,omitempty"`
	MaxAttempts *uint64                      `json:"maxAttempts,omitempty"`
	Window      *timeutil.Duration           `json:"window,omitempty"`
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
	Interval      timeutil.Duration `json:"interval,omitempty"`
	Timeout       timeutil.Duration `json:"timeout,omitempty"`
	StartPeriod   timeutil.Duration `json:"startPeriod,omitempty"`
	StartInterval timeutil.Duration `json:"startInterval,omitempty"`

	// Retries is the number of consecutive failures needed to consider a container as unhealthy.
	// Zero means inherit.
	Retries int `json:"retries,omitempty"`
}

type HealthcheckMode string

const (
	HealthcheckModeInherit  = HealthcheckMode("")
	HealthcheckModeNone     = HealthcheckMode("NONE")
	HealthcheckModeCmd      = HealthcheckMode("CMD")
	HealthcheckModeCmdShell = HealthcheckMode("CMD-SHELL")
)
