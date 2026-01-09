package docker

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

func ConvertFromServiceModeSpec(spec *swarm.ServiceSpec) *ServiceModeSpec {
	if spec == nil {
		return nil
	}
	res := &ServiceModeSpec{}
	switch {
	case spec.Mode.Replicated != nil:
		res.Mode = ServiceModeReplicated
		res.ServiceReplicas = spec.Mode.Replicated.Replicas
	case spec.Mode.ReplicatedJob != nil:
		res.Mode = ServiceModeReplicatedJob
		res.JobMaxConcurrent = spec.Mode.ReplicatedJob.MaxConcurrent
		res.JobTotalCompletions = spec.Mode.ReplicatedJob.TotalCompletions
	case spec.Mode.Global != nil:
		res.Mode = ServiceModeGlobal
	case spec.Mode.GlobalJob != nil:
		res.Mode = ServiceModeGlobalJob
	}
	return res
}

func ApplyServiceModeSpec(spec *swarm.ServiceSpec, req *ServiceModeSpec) {
	if req == nil || req.Mode == "" {
		return
	}
	replicated := spec.Mode.Replicated
	spec.Mode.Replicated = nil
	replicatedJob := spec.Mode.ReplicatedJob
	spec.Mode.ReplicatedJob = nil
	global := spec.Mode.Global
	spec.Mode.Global = nil
	globalJob := spec.Mode.GlobalJob
	spec.Mode.GlobalJob = nil

	switch req.Mode {
	case ServiceModeReplicated:
		if replicated == nil {
			replicated = &swarm.ReplicatedService{}
		}
		if req.ServiceReplicas != nil {
			replicated.Replicas = req.ServiceReplicas
		}
		spec.Mode.Replicated = replicated
	case ServiceModeReplicatedJob:
		if replicatedJob == nil {
			replicatedJob = &swarm.ReplicatedJob{}
		}
		if req.JobMaxConcurrent != nil {
			replicatedJob.MaxConcurrent = req.JobMaxConcurrent
		}
		if req.JobTotalCompletions != nil {
			replicatedJob.TotalCompletions = req.JobTotalCompletions
		}
		spec.Mode.ReplicatedJob = replicatedJob
	case ServiceModeGlobal:
		if global == nil {
			global = &swarm.GlobalService{}
		}
		spec.Mode.Global = global
	case ServiceModeGlobalJob:
		if globalJob == nil {
			globalJob = &swarm.GlobalJob{}
		}
		spec.Mode.GlobalJob = globalJob
	}
}

func ConvertFromServiceTaskSpec(taskSpec *swarm.TaskSpec) *TaskSpec {
	if taskSpec == nil {
		return nil
	}
	res := &TaskSpec{
		Networks:      ConvertFromServiceNetworks(taskSpec.Networks),
		Resources:     ConvertFromServiceResourceRequirements(taskSpec.Resources),
		Placement:     ConvertFromServicePlacement(taskSpec.Placement),
		RestartPolicy: ConvertFromServiceRestartPolicy(taskSpec.RestartPolicy),
	}
	return res
}

func ApplyServiceTaskSpec(taskSpec *swarm.TaskSpec, req *TaskSpec) {
	if req == nil {
		return
	}
	ApplyServiceNetworks(taskSpec, req.Networks)
	ApplyServiceResourceRequirements(taskSpec, req.Resources)
	ApplyServicePlacement(taskSpec, req.Placement)
	ApplyServiceRestartPolicy(taskSpec, req.RestartPolicy)
}

func ConvertFromServiceContainerSpec(contSpec *swarm.ContainerSpec) *ContainerSpec {
	if contSpec == nil {
		return nil
	}
	res := &ContainerSpec{
		Labels:           contSpec.Labels,
		Image:            &contSpec.Image,
		Command:          gofn.ToPtr(ConvertFromContainerCommand(contSpec.Command, contSpec.Args)),
		WorkingDir:       &contSpec.Dir,
		Hostname:         &contSpec.Hostname,
		User:             &contSpec.User,
		Groups:           contSpec.Groups,
		StopSignal:       &contSpec.StopSignal,
		TTY:              &contSpec.TTY,
		OpenStdin:        &contSpec.OpenStdin,
		ReadOnly:         &contSpec.ReadOnly,
		HostsFileEntries: ConvertFromContainerHosts(contSpec.Hosts),
		Sysctls:          contSpec.Sysctls,
		CapabilityAdd:    contSpec.CapabilityAdd,
		CapabilityDrop:   contSpec.CapabilityDrop,
		EnableGPU:        gofn.ToPtr(gofn.Contain(contSpec.CapabilityAdd, "[gpu]")),
		Ulimits:          ConvertFromServiceUlimits(contSpec.Ulimits),
		Healthcheck:      ConvertFromServiceHealthcheck(contSpec.Healthcheck),
	}
	if contSpec.StopGracePeriod != nil {
		res.StopGracePeriod = gofn.ToPtr(timeutil.Duration(*contSpec.StopGracePeriod))
	}
	res.BindMounts, res.VolumeMounts = ConvertFromServiceMounts(contSpec.Mounts)
	return res
}

func ApplyServiceContainerSpec(contSpec *swarm.ContainerSpec, req *ContainerSpec) {
	contSpec.Labels = req.Labels
	if req.Image != nil {
		contSpec.Image = *req.Image
	}
	if req.Command != nil {
		ApplyContainerCommand(contSpec, *req.Command)
	}
	if req.WorkingDir != nil {
		contSpec.Dir = *req.WorkingDir
	}
	if req.Hostname != nil {
		contSpec.Hostname = *req.Hostname
	}
	if req.User != nil {
		contSpec.User = *req.User
	}
	if req.Groups != nil {
		contSpec.Groups = req.Groups
	}
	if req.StopSignal != nil {
		contSpec.StopSignal = *req.StopSignal
	}
	if req.TTY != nil {
		contSpec.TTY = *req.TTY
	}
	if req.OpenStdin != nil {
		contSpec.OpenStdin = *req.OpenStdin
	}
	if req.ReadOnly != nil {
		contSpec.ReadOnly = *req.ReadOnly
	}
	if req.StopGracePeriod != nil {
		contSpec.StopGracePeriod = gofn.ToPtr(time.Duration(*req.StopGracePeriod))
	}

	ApplyContainerHosts(contSpec, req.HostsFileEntries)
	ApplyServiceMounts(contSpec, req.BindMounts, req.VolumeMounts)

	if req.Sysctls != nil {
		contSpec.Sysctls = req.Sysctls
	}
	if req.CapabilityAdd != nil {
		contSpec.CapabilityAdd = req.CapabilityAdd
	}
	if req.CapabilityDrop != nil {
		contSpec.CapabilityDrop = req.CapabilityDrop
	}
	if req.EnableGPU != nil {
		if *req.EnableGPU && !gofn.Contain(contSpec.CapabilityAdd, "[gpu]") {
			contSpec.CapabilityAdd = append(contSpec.CapabilityAdd, "[gpu]")
		} else {
			contSpec.CapabilityAdd = gofn.Drop(contSpec.CapabilityAdd, "[gpu]")
		}
	}

	ApplyServiceUlimits(contSpec, req.Ulimits)
	ApplyServiceHealthcheck(contSpec, req.Healthcheck)
}

func ConvertFromContainerCommand(cmd []string, args []string) string {
	return strings.Join(gofn.Concat(cmd, args), " ")
}

func ApplyContainerCommand(contSpec *swarm.ContainerSpec, cmd string) {
	contSpec.Command = gofn.StringSplit(cmd, " ", "\"")
}

func ConvertFromContainerHosts(hosts []string) (res []*HostsFileEntry) {
	for _, host := range hosts {
		parts := gofn.StringSplit(host, " ", "\"")
		res = append(res, &HostsFileEntry{
			Address:   parts[0],
			Hostnames: parts[1:],
		})
	}
	return res
}

func ApplyContainerHosts(contSpec *swarm.ContainerSpec, hosts []*HostsFileEntry) {
	contSpec.Hosts = make([]string, 0, len(hosts))
	for _, host := range hosts {
		s := append([]string{}, host.Address)
		s = append(s, host.Hostnames...)
		contSpec.Hosts = append(contSpec.Hosts, strings.Join(s, " "))
	}
}

func ConvertFromServiceUlimits(ulimits []*container.Ulimit) (res []*Ulimit) {
	for _, limit := range ulimits {
		if limit == nil {
			continue
		}
		res = append(res, &Ulimit{
			Name: limit.Name,
			Hard: limit.Hard,
			Soft: limit.Soft,
		})
	}
	return res
}

func ApplyServiceUlimits(contSpec *swarm.ContainerSpec, ulimits []*Ulimit) {
	contSpec.Ulimits = make([]*container.Ulimit, 0, len(ulimits))
	for _, limit := range ulimits {
		if limit == nil {
			continue
		}
		contSpec.Ulimits = append(contSpec.Ulimits, &container.Ulimit{
			Name: limit.Name,
			Hard: limit.Hard,
			Soft: limit.Soft,
		})
	}
}

func ConvertFromServiceMounts(mounts []mount.Mount) (bindMounts []*BindMount, volMounts []*VolumeMount) {
	for _, mnt := range mounts {
		switch mnt.Type { //nolint:exhaustive
		case mount.TypeBind:
			bindMounts = append(bindMounts, &BindMount{
				Source:   mnt.Source,
				Target:   mnt.Target,
				ReadOnly: mnt.ReadOnly,
			})
		case mount.TypeVolume:
			volMounts = append(volMounts, &VolumeMount{
				Source:   mnt.Source,
				Target:   mnt.Target,
				ReadOnly: mnt.ReadOnly,
			})
		}
	}
	return
}

func ApplyServiceMounts(contSpec *swarm.ContainerSpec, bindMounts []*BindMount, volMounts []*VolumeMount) {
	contSpec.Mounts = make([]mount.Mount, 0, len(bindMounts)+len(volMounts))
	for _, mnt := range bindMounts {
		contSpec.Mounts = append(contSpec.Mounts, mount.Mount{
			Type:     mount.TypeBind,
			Source:   mnt.Source,
			Target:   mnt.Target,
			ReadOnly: mnt.ReadOnly,
		})
	}
	for _, mnt := range volMounts {
		contSpec.Mounts = append(contSpec.Mounts, mount.Mount{
			Type:     mount.TypeVolume,
			Source:   mnt.Source,
			Target:   mnt.Target,
			ReadOnly: mnt.ReadOnly,
		})
	}
}

func ConvertFromServiceHealthcheck(config *container.HealthConfig) *Healthcheck {
	if config == nil {
		return nil
	}
	cmd := config.Test
	var mode HealthcheckMode
	if len(cmd) > 0 {
		mode = HealthcheckMode(cmd[0])
		cmd = cmd[1:]
	}
	res := &Healthcheck{
		Enabled:       mode != HealthcheckModeNone,
		Mode:          mode,
		Command:       strings.Join(cmd, " "),
		Interval:      timeutil.Duration(config.Interval),
		Timeout:       timeutil.Duration(config.Timeout),
		StartPeriod:   timeutil.Duration(config.StartPeriod),
		StartInterval: timeutil.Duration(config.StartInterval),
		Retries:       config.Retries,
	}
	return res
}

func ApplyServiceHealthcheck(contSpec *swarm.ContainerSpec, healthcheck *Healthcheck) {
	if healthcheck == nil {
		contSpec.Healthcheck = nil
		return
	}
	cmd := gofn.StringSplit(healthcheck.Command, " ", "\"")
	contSpec.Healthcheck = &container.HealthConfig{
		Test:          gofn.Concat([]string{string(healthcheck.Mode)}, cmd),
		Interval:      time.Duration(healthcheck.Interval),
		Timeout:       time.Duration(healthcheck.Timeout),
		StartPeriod:   time.Duration(healthcheck.StartPeriod),
		StartInterval: time.Duration(healthcheck.StartInterval),
		Retries:       healthcheck.Retries,
	}
}

func ConvertFromServiceNetworks(networks []swarm.NetworkAttachmentConfig) (res []*NetworkAttachment) {
	for _, net := range networks {
		res = append(res, &NetworkAttachment{
			Target:  net.Target,
			Aliases: net.Aliases,
		})
	}
	return res
}

func ApplyServiceNetworks(taskSpec *swarm.TaskSpec, networks []*NetworkAttachment) {
	taskSpec.Networks = make([]swarm.NetworkAttachmentConfig, 0, len(networks))
	for _, net := range networks {
		taskSpec.Networks = append(taskSpec.Networks, swarm.NetworkAttachmentConfig{
			Target:  net.Target,
			Aliases: net.Aliases,
		})
	}
}

func ConvertFromServiceResourceRequirements(req *swarm.ResourceRequirements) *ResourceRequirements {
	if req == nil {
		return nil
	}
	res := &ResourceRequirements{}

	if req.Reservations != nil {
		res.Reservations = &Resources{
			CPUs:     float64(req.Reservations.NanoCPUs / UnitCPUNano),
			MemoryMB: req.Reservations.MemoryBytes / UnitMemMB,
		}
		for _, r := range req.Reservations.GenericResources {
			if r.NamedResourceSpec != nil {
				res.Reservations.GenericResources = append(res.Reservations.GenericResources,
					fmt.Sprintf("%v=%v", r.NamedResourceSpec.Kind, r.NamedResourceSpec.Value))
			}
			if r.DiscreteResourceSpec != nil {
				res.Reservations.GenericResources = append(res.Reservations.GenericResources,
					fmt.Sprintf("%v=%v", r.DiscreteResourceSpec.Kind, r.DiscreteResourceSpec.Value))
			}
		}
	}

	if req.Limits != nil {
		res.Limits = &ResourceLimit{
			CPUs:     float64(req.Limits.NanoCPUs / UnitCPUNano),
			MemoryMB: req.Limits.MemoryBytes / UnitMemMB,
			Pids:     req.Limits.Pids,
		}
	}
	return res
}

func ApplyServiceResourceRequirements(taskSpec *swarm.TaskSpec, req *ResourceRequirements) {
	if req == nil {
		taskSpec.Resources = nil
		return
	}
	if taskSpec.Resources == nil {
		taskSpec.Resources = &swarm.ResourceRequirements{}
	}
	taskSpec.Resources.Reservations = nil
	taskSpec.Resources.Limits = nil

	if req.Reservations != nil {
		taskSpec.Resources.Reservations = &swarm.Resources{
			NanoCPUs:    int64(req.Reservations.CPUs * UnitCPUNano),
			MemoryBytes: req.Reservations.MemoryMB * UnitMemMB,
		}
		for _, r := range req.Reservations.GenericResources {
			k, v, _ := strings.Cut(r, "=")
			k, v = strings.TrimSpace(k), strings.TrimSpace(v)
			num, err := strconv.ParseInt(v, 10, 64)
			genericRes := swarm.GenericResource{}
			if err != nil {
				genericRes.NamedResourceSpec = &swarm.NamedGenericResource{
					Kind:  k,
					Value: v,
				}
			} else {
				genericRes.DiscreteResourceSpec = &swarm.DiscreteGenericResource{
					Kind:  k,
					Value: num,
				}
			}
			taskSpec.Resources.Reservations.GenericResources =
				append(taskSpec.Resources.Reservations.GenericResources, genericRes)
		}
	}

	if req.Limits != nil {
		taskSpec.Resources.Limits = &swarm.Limit{
			NanoCPUs:    int64(req.Limits.CPUs * UnitCPUNano),
			MemoryBytes: req.Limits.MemoryMB * UnitMemMB,
			Pids:        req.Limits.Pids,
		}
	}
}

func ConvertFromServicePlacement(placement *swarm.Placement) *Placement {
	if placement == nil {
		return nil
	}
	res := &Placement{
		Constraints: placement.Constraints,
	}
	for _, pref := range placement.Preferences {
		if pref.Spread != nil {
			res.Preferences = append(res.Preferences,
				fmt.Sprintf("%v=%v", "spread", pref.Spread.SpreadDescriptor))
		}
	}
	return res
}

func ApplyServicePlacement(taskSpec *swarm.TaskSpec, placement *Placement) {
	if placement == nil {
		taskSpec.Placement = nil
		return
	}
	if taskSpec.Placement == nil {
		taskSpec.Placement = &swarm.Placement{}
	}
	taskSpec.Placement.Constraints = placement.Constraints
	taskSpec.Placement.Preferences = nil

	for _, pref := range placement.Preferences {
		k, v, _ := strings.Cut(pref, "=")
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)

		if k == "spread" {
			taskSpec.Placement.Preferences = append(taskSpec.Placement.Preferences,
				swarm.PlacementPreference{
					Spread: &swarm.SpreadOver{SpreadDescriptor: v},
				})
		}
	}
}

func ConvertFromServiceRestartPolicy(policy *swarm.RestartPolicy) *RestartPolicy {
	if policy == nil {
		return nil
	}
	res := &RestartPolicy{
		Condition:   policy.Condition,
		MaxAttempts: policy.MaxAttempts,
	}
	if policy.Delay != nil {
		res.Delay = gofn.ToPtr(timeutil.Duration(*policy.Delay))
	}
	if policy.Window != nil {
		res.Window = gofn.ToPtr(timeutil.Duration(*policy.Window))
	}
	return res
}

func ApplyServiceRestartPolicy(taskSpec *swarm.TaskSpec, policy *RestartPolicy) {
	if policy == nil {
		taskSpec.RestartPolicy = nil
		return
	}
	if taskSpec.RestartPolicy == nil {
		taskSpec.RestartPolicy = &swarm.RestartPolicy{}
	}
	taskSpec.RestartPolicy.Condition = policy.Condition
	taskSpec.RestartPolicy.MaxAttempts = policy.MaxAttempts
	taskSpec.RestartPolicy.Delay = nil
	taskSpec.RestartPolicy.Window = nil
	if policy.Delay != nil {
		taskSpec.RestartPolicy.Delay = gofn.ToPtr(time.Duration(*policy.Delay))
	}
	if policy.Window != nil {
		taskSpec.RestartPolicy.Window = gofn.ToPtr(time.Duration(*policy.Window))
	}
}

func ConvertFromServiceEndpointSpec(endpointSpec *swarm.EndpointSpec) *EndpointSpec {
	if endpointSpec == nil {
		return nil
	}
	res := &EndpointSpec{
		Mode: endpointSpec.Mode,
	}
	for _, port := range endpointSpec.Ports {
		res.Ports = append(res.Ports, &PortConfig{
			Target:      port.TargetPort,
			Published:   port.PublishedPort,
			Protocol:    port.Protocol,
			PublishMode: port.PublishMode,
		})
	}
	return res
}

func ApplyServiceEndpointSpec(spec *swarm.ServiceSpec, endpointSpec *EndpointSpec) {
	if endpointSpec == nil {
		spec.EndpointSpec = nil
		return
	}
	if spec.EndpointSpec == nil {
		spec.EndpointSpec = &swarm.EndpointSpec{}
	}
	spec.EndpointSpec.Mode = endpointSpec.Mode
	spec.EndpointSpec.Ports = make([]swarm.PortConfig, 0, len(endpointSpec.Ports))
	for _, port := range endpointSpec.Ports {
		spec.EndpointSpec.Ports = append(spec.EndpointSpec.Ports, swarm.PortConfig{
			TargetPort:    port.Target,
			PublishedPort: port.Published,
			Protocol:      port.Protocol,
			PublishMode:   port.PublishMode,
		})
	}
}
