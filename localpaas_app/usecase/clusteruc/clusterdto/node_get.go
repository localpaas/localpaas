package clusterdto

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	nodeIDMaxLen = 100
)

type GetNodeReq struct {
	NodeID string `json:"-"`
}

func NewGetNodeReq() *GetNodeReq {
	return &GetNodeReq{}
}

func (req *GetNodeReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	// NOTE: node id is docker id, it's not ULID
	validators = append(validators, basedto.ValidateStr(&req.NodeID, true, 1, nodeIDMaxLen, "nodeId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetNodeResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *NodeResp         `json:"data"`
}

type NodeResp struct {
	ID        string            `json:"id"`
	Hostname  string            `json:"hostname"`
	Addr      string            `json:"addr"`
	Status    base.NodeStatus   `json:"status"`
	Role      base.NodeRole     `json:"role"`
	IsLeader  bool              `json:"isLeader"`
	Platform  *NodePlatformResp `json:"platform"`
	Resources *NodeResources    `json:"resources"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NodePlatformResp struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

type NodeResources struct {
	NanoCPUs    int64 `json:"nanoCPUs"`
	MemoryBytes int64 `json:"memoryBytes"`
}

func TransformNode(node *swarm.Node) *NodeResp {
	isManager := node.Spec.Role == swarm.NodeRoleManager
	return &NodeResp{
		ID:       node.ID,
		Status:   base.NodeStatus(node.Status.State),
		Role:     base.NodeRole(node.Spec.Role),
		IsLeader: isManager && node.ManagerStatus != nil && node.ManagerStatus.Leader,
		Hostname: node.Description.Hostname,
		Addr:     node.Status.Addr,
		Platform: &NodePlatformResp{
			Architecture: node.Description.Platform.Architecture,
			OS:           node.Description.Platform.OS,
		},
		Resources: &NodeResources{
			NanoCPUs:    node.Description.Resources.NanoCPUs,
			MemoryBytes: node.Description.Resources.MemoryBytes,
		},
		CreatedAt: node.CreatedAt,
		UpdatedAt: node.UpdatedAt,
	}
}
