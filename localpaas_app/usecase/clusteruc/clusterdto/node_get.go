package clusterdto

import (
	"time"

	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type GetNodeReq struct {
	NodeID string `json:"-"`
}

func NewGetNodeReq() *GetNodeReq {
	return &GetNodeReq{}
}

func (req *GetNodeReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.NodeID, true, "nodeId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetNodeResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *NodeResp         `json:"data"`
}

type NodeResp struct {
	ID          string          `json:"id"`
	HostName    string          `json:"hostName"`
	IP          string          `json:"ip"`
	Status      base.NodeStatus `json:"status"`
	InfraStatus string          `json:"infraStatus"`
	IsLeader    bool            `json:"isLeader"`
	IsManager   bool            `json:"isManager"`

	LastSyncedAt time.Time `json:"lastSyncedAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type NodeBaseResp struct {
	ID           string          `json:"id"`
	HostName     string          `json:"hostName"`
	IP           string          `json:"ip"`
	Status       base.NodeStatus `json:"status"`
	InfraStatus  string          `json:"infraStatus"`
	IsLeader     bool            `json:"isLeader"`
	IsManager    bool            `json:"isManager"`
	LastSyncedAt time.Time       `json:"lastSyncedAt"`
}

func TransformNode(node *entity.Node) (resp *NodeResp, err error) {
	if err = copier.Copy(&resp, &node); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

func TransformNodesBase(nodes []*entity.Node) []*NodeBaseResp {
	return gofn.MapSlice(nodes, func(node *entity.Node) *NodeBaseResp {
		return &NodeBaseResp{
			ID:          node.ID,
			HostName:    node.HostName,
			IP:          node.IP,
			Status:      node.Status,
			InfraStatus: node.InfraStatus,
			IsLeader:    node.IsLeader,
			IsManager:   node.IsManager,
		}
	})
}
