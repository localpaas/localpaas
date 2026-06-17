package appsettingsdto

import (
	"sort"
	"time"

	"github.com/moby/moby/api/types/swarm"
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/nodeuc/nodedto"
)

type GetAppServiceTasksReq struct {
	ProjectID string            `json:"-"`
	AppID     string            `json:"-"`
	States    []swarm.TaskState `json:"-" mapstructure:"state"`
}

func NewGetAppServiceTasksReq() *GetAppServiceTasksReq {
	return &GetAppServiceTasksReq{
		States: []swarm.TaskState{swarm.TaskStateRunning, swarm.TaskStateComplete, swarm.TaskStateShutdown,
			swarm.TaskStateFailed},
	}
}

func (req *GetAppServiceTasksReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAppServiceTasksResp struct {
	Meta *basedto.Meta      `json:"meta"`
	Data []*ServiceTaskResp `json:"data"`
}

type ServiceTaskResp struct {
	ID           string                 `json:"id"`
	Slot         int                    `json:"slot"`
	Node         *nodedto.NodeBaseResp  `json:"node"`
	Status       *ServiceTaskStatusResp `json:"status"`
	DesiredState swarm.TaskState        `json:"desiredState"`
}

type ServiceTaskStatusResp struct {
	Timestamp       time.Time        `json:"timestamp"`
	State           swarm.TaskState  `json:"state"`
	Message         string           `json:"message"`
	Err             string           `json:"err"`
	ContainerStatus *ContainerStatus `json:"containerStatus"`
}

type ContainerStatus struct {
	ContainerID string `json:"containerId"`
}

func TransformServiceTask(task *swarm.Task, nodeMap map[string]*swarm.Node) (resp *ServiceTaskResp, err error) {
	if err = copier.Copy(&resp, task); err != nil {
		return nil, apperrors.Wrap(err)
	}

	// TODO: remove this later
	resp.DesiredState = "[" + resp.DesiredState + "]"

	// Transform node details
	resp.Node = nodedto.TransformNodeBase(nodeMap[task.NodeID])

	return resp, nil
}

func TransformServiceTasks(tasks []swarm.Task, nodes []swarm.Node) (resp []*ServiceTaskResp, err error) {
	// Build node map
	nodeMap := make(map[string]*swarm.Node, len(nodes))
	for _, node := range nodes {
		nodeMap[node.ID] = &node
	}

	resp = make([]*ServiceTaskResp, 0, len(tasks))
	for i := range tasks {
		taskResp, err := TransformServiceTask(&tasks[i], nodeMap)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, taskResp)
	}

	//nolint:mnd,exhaustive
	stateRank := func(state swarm.TaskState) int {
		switch state {
		case swarm.TaskStateRunning:
			return 0
		case swarm.TaskStateComplete:
			return 1
		case swarm.TaskStateShutdown:
			return 2
		case swarm.TaskStateFailed:
			return 3
		default:
			return 4
		}
	}

	sort.Slice(resp, func(i, j int) bool {
		var stateI, stateJ swarm.TaskState
		var slotI, slotJ int
		var timeI, timeJ time.Time

		if resp[i].Status != nil {
			stateI = resp[i].Status.State
			slotI = resp[i].Slot
			timeI = resp[i].Status.Timestamp
		}
		if resp[j].Status != nil {
			stateJ = resp[j].Status.State
			slotJ = resp[j].Slot
			timeJ = resp[j].Status.Timestamp
		}

		rankI := stateRank(stateI)
		rankJ := stateRank(stateJ)
		if rankI != rankJ {
			return rankI < rankJ
		}
		if slotI != slotJ {
			return slotI < slotJ
		}
		return timeI.After(timeJ)
	})

	return resp, nil
}
