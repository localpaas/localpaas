package sshkeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc/sshkeydto"
)

func (uc *SSHKeyUC) GetSSHKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *sshkeydto.GetSSHKeyReq,
) (*sshkeydto.GetSSHKeyResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, req.ID,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeSSHKey),
		bunex.SelectRelation("ObjectAccesses",
			bunex.SelectWhere("acl_permission.subject_type IN (?)", bunex.In([]base.SubjectType{
				base.SubjectTypeProject, base.SubjectTypeApp,
			})),
			bunex.SelectRelation("SubjectProject"),
			bunex.SelectRelation("SubjectApp"),
		),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := sshkeydto.TransformSSHKey(setting, true)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sshkeydto.GetSSHKeyResp{
		Data: resp,
	}, nil
}
