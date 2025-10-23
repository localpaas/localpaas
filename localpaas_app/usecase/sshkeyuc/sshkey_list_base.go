package sshkeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sshkeyuc/sshkeydto"
)

func (uc *SSHKeyUC) ListSSHKeyBase(
	ctx context.Context,
	auth *basedto.Auth,
	req *sshkeydto.ListSSHKeyBaseReq,
) (*sshkeydto.ListSSHKeyBaseResp, error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.type = ?", base.SettingTypeSSHKey),
	}

	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhereGroup(
				bunex.SelectWhere("setting.name ILIKE ?", keyword),
			),
		)
	}

	settings, pagingMeta, err := uc.settingRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := sshkeydto.TransformSSHKeysBase(settings)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sshkeydto.ListSSHKeyBaseResp{
		Meta: &basedto.Meta{Page: pagingMeta},
		Data: resp,
	}, nil
}
