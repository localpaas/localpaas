package clusteruc

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc/clusterdto"
	"github.com/localpaas/localpaas/services/ssh"
)

const (
	executionTimeout = time.Second * 15
)

func (uc *ClusterUC) JoinNode(
	ctx context.Context,
	auth *basedto.Auth,
	req *clusterdto.JoinNodeReq,
) (*clusterdto.JoinNodeResp, error) {
	data := &joinNodeData{}
	err := uc.loadJoinNodeData(ctx, uc.db, req, data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	cmdCtx, cancelFunc := context.WithTimeout(ctx, executionTimeout)
	defer cancelFunc()

	command := fmt.Sprintf("docker swarm leave --force && docker swarm join --token %s %s",
		data.JoinToken, data.PreferManagerAddr)
	output, err := ssh.Execute(cmdCtx, &ssh.CommandInput{
		Host:       req.Host,
		Port:       req.Port,
		User:       req.User,
		PrivateKey: data.SSHKey.PrivateKey,
		Passphrase: data.SSHKey.Passphrase,
		Command:    command,
	})

	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}

	return &clusterdto.JoinNodeResp{
		Data: &clusterdto.JoinNodeDataResp{
			Success:       err == nil,
			ErrorMessage:  errorMessage,
			CommandOutput: output,
		},
	}, nil
}

type joinNodeData struct {
	SSHKey            *entity.SSHKey
	JoinToken         string
	PreferManagerAddr string
}

func (uc *ClusterUC) loadJoinNodeData(
	ctx context.Context,
	db database.IDB,
	req *clusterdto.JoinNodeReq,
	data *joinNodeData,
) error {
	sshKeySetting, err := uc.settingRepo.GetByID(ctx, db, req.SSHKey.ID,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeSSHKey),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}

	sshKey, err := sshKeySetting.ParseSSHKey(true)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.SSHKey = sshKey

	// Find join token from the cluster
	theSwarm, err := uc.dockerManager.SwarmInspect(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}

	joinToken := gofn.If(req.JoinAsManager, theSwarm.JoinTokens.Manager, theSwarm.JoinTokens.Worker) //nolint
	if joinToken == "" {
		return apperrors.Wrap(apperrors.ErrDockerJoinTokenNotFound)
	}
	data.JoinToken = joinToken

	// List all manager nodes to get the addr to join new node
	managerNodes, err := uc.dockerManager.NodeList(ctx, func(opts *swarm.NodeListOptions) {
		opts.Filters = filters.NewArgs(filters.Arg("role", "manager"))
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	var leaderAddr, managerAddr string
	for _, node := range managerNodes {
		mgrStatus := node.ManagerStatus
		if mgrStatus.Reachability == swarm.ReachabilityReachable {
			managerAddr = mgrStatus.Addr
			if mgrStatus.Leader {
				leaderAddr = mgrStatus.Addr
			}
		}
	}
	data.PreferManagerAddr = gofn.Coalesce(leaderAddr, managerAddr)
	if data.PreferManagerAddr == "" {
		return apperrors.Wrap(apperrors.ErrDockerActiveManagerNodeNotFound)
	}

	return nil
}
