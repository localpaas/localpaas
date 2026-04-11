package docker

import (
	"strings"

	"github.com/docker/docker/api/types/swarm"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/pkg/shellutil"
)

const (
	serviceRestrictedLabelsPrefix = "localpaas."
)

var (
	serviceSystemLabels = map[string]struct{}{
		StackLabelNamespace: {},
	}
)

func ContainerCommandBuild(cmd []string, args []string) string {
	return strings.Join(gofn.Concat(cmd, args), " ")
}

func ContainerCommandApply(contSpec *swarm.ContainerSpec, cmd string) {
	if cmd == "" {
		contSpec.Command = nil
	} else {
		contSpec.Command = gofn.Must(shellutil.CmdSplit(cmd))
	}
}

func ServiceFilterOutSystemLabels(labels map[string]string) map[string]string {
	resp := make(map[string]string, len(labels))
	for k, v := range labels {
		if _, exists := serviceSystemLabels[k]; exists {
			continue
		}
		if strings.HasPrefix(k, serviceRestrictedLabelsPrefix) {
			continue
		}
		resp[k] = v
	}
	return resp
}

func ServiceValidateUserLabels(labels map[string]string, stopAtFirstViolation bool) (unallowedLabels []string) {
	for k := range labels {
		if _, exists := serviceSystemLabels[k]; exists {
			if stopAtFirstViolation {
				return []string{k}
			}
			unallowedLabels = append(unallowedLabels, k)
			continue
		}
		if strings.HasPrefix(k, serviceRestrictedLabelsPrefix) {
			if stopAtFirstViolation {
				return []string{k}
			}
			unallowedLabels = append(unallowedLabels, k)
			continue
		}
	}
	return unallowedLabels
}

func ServiceApplyUserLabels(currLabels, userLabels map[string]string) map[string]string {
	appliedLabels := make(map[string]string, len(userLabels))
	for k, v := range currLabels {
		if _, exists := serviceSystemLabels[k]; exists {
			appliedLabels[k] = v
			continue
		}
		if strings.HasPrefix(k, serviceRestrictedLabelsPrefix) {
			appliedLabels[k] = v
			continue
		}
	}
	for k, v := range userLabels {
		if _, exists := serviceSystemLabels[k]; exists {
			continue
		}
		if strings.HasPrefix(k, serviceRestrictedLabelsPrefix) {
			continue
		}
		appliedLabels[k] = v
	}
	return appliedLabels
}
