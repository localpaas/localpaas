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
	if req == nil {
		return
	}
	currMode := &spec.Mode
	spec.Mode = swarm.ServiceMode{}
	switch req.Mode {
	case ServiceModeReplicated:
		item := currMode.Replicated
		if item == nil {
			item = &swarm.ReplicatedService{}
		}
		item.Replicas = req.ServiceReplicas
		spec.Mode.Replicated = item
	case ServiceModeReplicatedJob:
		item := currMode.ReplicatedJob
		if item == nil {
			item = &swarm.ReplicatedJob{}
		}
		item.MaxConcurrent = req.JobMaxConcurrent
		item.TotalCompletions = req.JobTotalCompletions
		spec.Mode.ReplicatedJob = item
	case ServiceModeGlobal:
		item := currMode.Global
		if item == nil {
			item = &swarm.GlobalService{}
		}
		spec.Mode.Global = item
	case ServiceModeGlobalJob:
		item := currMode.GlobalJob
		if item == nil {
			item = &swarm.GlobalJob{}
		}
		spec.Mode.GlobalJob = item
	}
}

func ConvertFromServiceTaskSpec(taskSpec *swarm.TaskSpec) *TaskSpec {
	if taskSpec == nil {
		return nil
	}
	res := &TaskSpec{
		Networks:      ConvertFromServiceNetworks(taskSpec.Networks),
		Resources:     ConvertFromServiceResources(taskSpec.Resources),
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
	ApplyServiceResources(taskSpec, req.Resources)
	ApplyServicePlacement(taskSpec, req.Placement)
	ApplyServiceRestartPolicy(taskSpec, req.RestartPolicy)
}

func ConvertFromServiceContainerSpec(contSpec *swarm.ContainerSpec) *ContainerSpec {
	if contSpec == nil {
		return nil
	}
	res := &ContainerSpec{
		Labels:           contSpec.Labels,
		Image:            contSpec.Image,
		Command:          ConvertFromServiceCommand(contSpec.Command, contSpec.Args),
		WorkingDir:       contSpec.Dir,
		Hostname:         contSpec.Hostname,
		User:             contSpec.User,
		Groups:           contSpec.Groups,
		StopSignal:       contSpec.StopSignal,
		TTY:              contSpec.TTY,
		OpenStdin:        contSpec.OpenStdin,
		ReadOnly:         contSpec.ReadOnly,
		Privileges:       ConvertFromServicePrivileges(contSpec.Privileges),
		HostsFileEntries: ConvertFromServiceHosts(contSpec.Hosts),
		DNSConfig:        ConvertFromServiceDNSConfig(contSpec.DNSConfig),
		Sysctls:          contSpec.Sysctls,
		CapabilityAdd:    contSpec.CapabilityAdd,
		CapabilityDrop:   contSpec.CapabilityDrop,
		EnableGPU:        gofn.Contain(contSpec.CapabilityAdd, "[gpu]"),
		Ulimits:          ConvertFromServiceUlimits(contSpec.Ulimits),
		Healthcheck:      ConvertFromServiceHealthcheck(contSpec.Healthcheck),
		Mounts:           ConvertFromServiceMounts(contSpec.Mounts),
	}
	if contSpec.StopGracePeriod != nil {
		res.StopGracePeriod = gofn.ToPtr(timeutil.Duration(*contSpec.StopGracePeriod))
	}
	return res
}

func ApplyServiceContainerSpec(contSpec *swarm.ContainerSpec, req *ContainerSpec) {
	contSpec.Labels = req.Labels
	contSpec.Image = req.Image
	ApplyServiceCommand(contSpec, req.Command)
	contSpec.Dir = req.WorkingDir
	contSpec.Hostname = req.Hostname
	contSpec.User = req.User
	contSpec.Groups = req.Groups
	contSpec.StopSignal = req.StopSignal
	contSpec.TTY = req.TTY
	contSpec.OpenStdin = req.OpenStdin
	contSpec.ReadOnly = req.ReadOnly
	if req.StopGracePeriod != nil {
		contSpec.StopGracePeriod = gofn.ToPtr(time.Duration(*req.StopGracePeriod))
	}

	ApplyServicePrivileges(contSpec, req.Privileges)
	ApplyServiceHosts(contSpec, req.HostsFileEntries)
	ApplyServiceDNSConfig(contSpec, req.DNSConfig)
	ApplyServiceMounts(contSpec, req.Mounts)

	contSpec.Sysctls = req.Sysctls
	contSpec.CapabilityAdd = req.CapabilityAdd
	contSpec.CapabilityDrop = req.CapabilityDrop
	if req.EnableGPU && !gofn.Contain(contSpec.CapabilityAdd, "[gpu]") {
		contSpec.CapabilityAdd = append(contSpec.CapabilityAdd, "[gpu]")
	} else if !req.EnableGPU {
		contSpec.CapabilityAdd = gofn.Drop(contSpec.CapabilityAdd, "[gpu]")
	}

	ApplyServiceUlimits(contSpec, req.Ulimits)
	ApplyServiceHealthcheck(contSpec, req.Healthcheck)
}

func ConvertFromServiceCommand(cmd []string, args []string) string {
	return strings.Join(gofn.Concat(cmd, args), " ")
}

func ApplyServiceCommand(contSpec *swarm.ContainerSpec, cmd string) {
	contSpec.Command = gofn.StringSplit(cmd, " ", "\"")
}

func ConvertFromServicePrivileges(privileges *swarm.Privileges) (res *Privileges) {
	if privileges == nil {
		return nil
	}
	res = &Privileges{
		NoNewPrivileges: privileges.NoNewPrivileges,
	}
	if privileges.SELinuxContext != nil {
		res.SELinuxContext = &SELinuxContext{
			Disable: privileges.SELinuxContext.Disable,
			User:    privileges.SELinuxContext.User,
			Role:    privileges.SELinuxContext.Role,
			Type:    privileges.SELinuxContext.Type,
			Level:   privileges.SELinuxContext.Level,
		}
	}
	if privileges.Seccomp != nil {
		res.Seccomp = &SeccompOpts{
			Mode:    privileges.Seccomp.Mode,
			Profile: privileges.Seccomp.Profile,
		}
	}
	if privileges.AppArmor != nil {
		res.AppArmor = &AppArmorOpts{
			Mode: privileges.AppArmor.Mode,
		}
	}
	return res
}

func ApplyServicePrivileges(contSpec *swarm.ContainerSpec, privileges *Privileges) {
	if privileges == nil {
		return
	}
	if contSpec.Privileges == nil {
		contSpec.Privileges = &swarm.Privileges{}
	}
	contSpec.Privileges.NoNewPrivileges = privileges.NoNewPrivileges

	if privileges.SELinuxContext != nil {
		if contSpec.Privileges.SELinuxContext == nil {
			contSpec.Privileges.SELinuxContext = &swarm.SELinuxContext{}
		}
		contSpec.Privileges.SELinuxContext.Disable = privileges.SELinuxContext.Disable
		contSpec.Privileges.SELinuxContext.User = privileges.SELinuxContext.User
		contSpec.Privileges.SELinuxContext.Role = privileges.SELinuxContext.Role
		contSpec.Privileges.SELinuxContext.Type = privileges.SELinuxContext.Type
		contSpec.Privileges.SELinuxContext.Level = privileges.SELinuxContext.Level
	} else {
		contSpec.Privileges.SELinuxContext = nil
	}

	if privileges.Seccomp != nil {
		if contSpec.Privileges.Seccomp == nil {
			contSpec.Privileges.Seccomp = &swarm.SeccompOpts{}
		}
		contSpec.Privileges.Seccomp.Mode = privileges.Seccomp.Mode
		contSpec.Privileges.Seccomp.Profile = privileges.Seccomp.Profile
	} else {
		contSpec.Privileges.Seccomp = nil
	}

	if privileges.AppArmor != nil {
		if contSpec.Privileges.AppArmor == nil {
			contSpec.Privileges.AppArmor = &swarm.AppArmorOpts{}
		}
		contSpec.Privileges.AppArmor.Mode = privileges.AppArmor.Mode
	} else {
		contSpec.Privileges.AppArmor = nil
	}
}

func ConvertFromServiceHosts(hosts []string) (res []*HostsFileEntry) {
	res = make([]*HostsFileEntry, 0, len(hosts))
	for _, host := range hosts {
		parts := gofn.StringSplit(host, " ", "\"")
		res = append(res, &HostsFileEntry{
			Address:   parts[0],
			Hostnames: parts[1:],
		})
	}
	return res
}

func ApplyServiceHosts(contSpec *swarm.ContainerSpec, hosts []*HostsFileEntry) {
	contSpec.Hosts = make([]string, 0, len(hosts))
	for _, host := range hosts {
		s := append([]string{}, host.Address)
		s = append(s, host.Hostnames...)
		contSpec.Hosts = append(contSpec.Hosts, strings.Join(s, " "))
	}
}

func ConvertFromServiceDNSConfig(config *swarm.DNSConfig) (res *DNSConfig) {
	if config == nil {
		return nil
	}
	return &DNSConfig{
		Nameservers: config.Nameservers,
		Search:      config.Search,
		Options:     config.Options,
	}
}

func ApplyServiceDNSConfig(contSpec *swarm.ContainerSpec, config *DNSConfig) {
	if config == nil {
		return
	}
	if contSpec.DNSConfig == nil {
		contSpec.DNSConfig = &swarm.DNSConfig{}
	}
	contSpec.DNSConfig.Nameservers = config.Nameservers
	contSpec.DNSConfig.Search = config.Search
	contSpec.DNSConfig.Options = config.Options
}

func ConvertFromServiceUlimits(ulimits []*container.Ulimit) (res []*Ulimit) {
	res = make([]*Ulimit, 0, len(ulimits))
	for i, limit := range ulimits {
		if limit == nil {
			continue
		}
		res = append(res, &Ulimit{
			Index: gofn.ToPtr(i),
			Name:  limit.Name,
			Hard:  limit.Hard,
			Soft:  limit.Soft,
		})
	}
	return res
}

func ApplyServiceUlimits(contSpec *swarm.ContainerSpec, ulimits []*Ulimit) {
	currUlimits := contSpec.Ulimits
	contSpec.Ulimits = make([]*container.Ulimit, 0, len(ulimits))
	for _, limit := range ulimits {
		if limit == nil {
			continue
		}
		var item *container.Ulimit
		if limit.Index != nil {
			item = currUlimits[*limit.Index]
		} else {
			item = &container.Ulimit{}
		}
		item.Name = limit.Name
		item.Hard = limit.Hard
		item.Soft = limit.Soft
		contSpec.Ulimits = append(contSpec.Ulimits, item)
	}
}

func ConvertFromServiceMounts(mounts []mount.Mount) (res []*Mount) {
	res = make([]*Mount, 0, len(mounts))
	for i, mnt := range mounts {
		res = append(res, &Mount{
			Index:    gofn.ToPtr(i),
			Type:     mnt.Type,
			Source:   mnt.Source,
			Target:   mnt.Target,
			ReadOnly: mnt.ReadOnly,
		})
	}
	return
}

func ApplyServiceMounts(contSpec *swarm.ContainerSpec, mounts []*Mount) {
	currMounts := contSpec.Mounts
	contSpec.Mounts = make([]mount.Mount, 0, len(mounts))
	for _, mnt := range mounts {
		var item mount.Mount
		if mnt.Index != nil {
			item = currMounts[*mnt.Index]
		} else {
			item = mount.Mount{}
		}
		item.Type = mnt.Type
		item.Source = mnt.Source
		item.Target = mnt.Target
		item.ReadOnly = mnt.ReadOnly
		contSpec.Mounts = append(contSpec.Mounts, item)
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
		return
	}
	if contSpec.Healthcheck == nil {
		contSpec.Healthcheck = &container.HealthConfig{}
	}
	cmd := gofn.StringSplit(healthcheck.Command, " ", "\"")
	contSpec.Healthcheck.Test = gofn.Concat([]string{string(healthcheck.Mode)}, cmd)
	contSpec.Healthcheck.Interval = time.Duration(healthcheck.Interval)
	contSpec.Healthcheck.Timeout = time.Duration(healthcheck.Timeout)
	contSpec.Healthcheck.StartPeriod = time.Duration(healthcheck.StartPeriod)
	contSpec.Healthcheck.StartInterval = time.Duration(healthcheck.StartInterval)
	contSpec.Healthcheck.Retries = healthcheck.Retries
}

func ConvertFromServiceNetworks(networks []swarm.NetworkAttachmentConfig) (res []*NetworkAttachment) {
	res = make([]*NetworkAttachment, 0, len(networks))
	for i, net := range networks {
		res = append(res, &NetworkAttachment{
			Index:   gofn.ToPtr(i),
			Target:  net.Target,
			Aliases: net.Aliases,
		})
	}
	return res
}

func ApplyServiceNetworks(taskSpec *swarm.TaskSpec, networks []*NetworkAttachment) {
	currNetworks := taskSpec.Networks
	taskSpec.Networks = make([]swarm.NetworkAttachmentConfig, 0, len(networks))
	for _, net := range networks {
		var item swarm.NetworkAttachmentConfig
		if net.Index != nil {
			item = currNetworks[*net.Index]
		} else {
			item = swarm.NetworkAttachmentConfig{}
		}
		item.Target = net.Target
		item.Aliases = net.Aliases
		taskSpec.Networks = append(taskSpec.Networks, item)
	}
}

func ConvertFromServiceResources(req *swarm.ResourceRequirements) *ResourceRequirements {
	if req == nil {
		return nil
	}
	res := &ResourceRequirements{}

	if req.Reservations != nil {
		res.Reservations = &Resources{
			CPUs:             float64(req.Reservations.NanoCPUs / UnitCPUNano),
			MemoryMB:         req.Reservations.MemoryBytes / UnitMemMB,
			GenericResources: make([]string, 0, len(req.Reservations.GenericResources)),
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

func ApplyServiceResources(taskSpec *swarm.TaskSpec, req *ResourceRequirements) {
	if req == nil {
		return
	}
	if taskSpec.Resources == nil {
		taskSpec.Resources = &swarm.ResourceRequirements{}
	}

	if req.Reservations != nil {
		if taskSpec.Resources.Reservations == nil {
			taskSpec.Resources.Reservations = &swarm.Resources{}
		}
		taskSpec.Resources.Reservations.NanoCPUs = int64(req.Reservations.CPUs * UnitCPUNano)
		taskSpec.Resources.Reservations.MemoryBytes = req.Reservations.MemoryMB * UnitMemMB

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
		if taskSpec.Resources.Limits == nil {
			taskSpec.Resources.Limits = &swarm.Limit{}
		}
		taskSpec.Resources.Limits.NanoCPUs = int64(req.Limits.CPUs * UnitCPUNano)
		taskSpec.Resources.Limits.MemoryBytes = req.Limits.MemoryMB * UnitMemMB
		taskSpec.Resources.Limits.Pids = req.Limits.Pids
	}
}

func ConvertFromServicePlacement(placement *swarm.Placement) *Placement {
	if placement == nil {
		return nil
	}
	res := &Placement{
		Constraints: placement.Constraints,
		Preferences: make([]string, 0, len(placement.Preferences)),
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
		return
	}
	if taskSpec.Placement == nil {
		taskSpec.Placement = &swarm.Placement{}
	}
	taskSpec.Placement.Constraints = placement.Constraints
	taskSpec.Placement.Preferences = make([]swarm.PlacementPreference, 0, len(placement.Preferences))

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
		return
	}
	if taskSpec.RestartPolicy == nil {
		taskSpec.RestartPolicy = &swarm.RestartPolicy{}
	}
	taskSpec.RestartPolicy.Condition = policy.Condition
	taskSpec.RestartPolicy.MaxAttempts = policy.MaxAttempts
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
		Mode:  endpointSpec.Mode,
		Ports: make([]*PortConfig, 0, len(endpointSpec.Ports)),
	}
	for i, port := range endpointSpec.Ports {
		res.Ports = append(res.Ports, &PortConfig{
			Index:       gofn.ToPtr(i),
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
		return
	}
	if spec.EndpointSpec == nil {
		spec.EndpointSpec = &swarm.EndpointSpec{}
	}
	spec.EndpointSpec.Mode = endpointSpec.Mode

	currPorts := spec.EndpointSpec.Ports
	spec.EndpointSpec.Ports = make([]swarm.PortConfig, 0, len(endpointSpec.Ports))
	for _, port := range endpointSpec.Ports {
		var item swarm.PortConfig
		if port.Index != nil {
			item = currPorts[*port.Index]
		} else {
			item = swarm.PortConfig{}
		}
		item.TargetPort = port.Target
		item.PublishedPort = port.Published
		item.Protocol = port.Protocol
		item.PublishMode = port.PublishMode
		spec.EndpointSpec.Ports = append(spec.EndpointSpec.Ports, item)
	}
}
