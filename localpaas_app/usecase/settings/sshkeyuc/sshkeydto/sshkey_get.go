package sshkeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

const (
	maskedSecretKey = "****************"
)

type GetSSHKeyReq struct {
	ID string `json:"-"`
}

func NewGetSSHKeyReq() *GetSSHKeyReq {
	return &GetSSHKeyReq{}
}

func (req *GetSSHKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetSSHKeyResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *SSHKeyResp       `json:"data"`
}

type SSHKeyResp struct {
	ID              string                     `json:"id"`
	Name            string                     `json:"name"`
	PrivateKey      string                     `json:"privateKey"`
	Passphrase      string                     `json:"passphrase,omitempty"`
	Encrypted       bool                       `json:"encrypted,omitempty"`
	ProjectAccesses []*SSHKeyProjectAccessResp `json:"projectAccesses"`
}

type SSHKeyProjectAccessResp struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Allowed     bool                   `json:"allowed"`
	AppAccesses []*SSHKeyAppAccessResp `json:"appAccesses"`
}

type SSHKeyAppAccessResp struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Allowed bool   `json:"allowed"`
}

func TransformSSHKey(setting *entity.Setting, decrypt bool) (resp *SSHKeyResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	sshKey, err := setting.ParseSSHKey(decrypt)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, &sshKey); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = sshKey.IsEncrypted()
	if resp.Encrypted {
		resp.PrivateKey = maskedSecretKey
		resp.Passphrase = maskedSecretKey
	}

	resp.ProjectAccesses, err = TransformSSHKeyObjectAccesses(setting.ObjectAccesses)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

func TransformSSHKeyObjectAccesses(accesses []*entity.ACLPermission) (
	resp []*SSHKeyProjectAccessResp, err error) {
	mapProjectResp := make(map[string]*SSHKeyProjectAccessResp)
	for _, access := range accesses {
		if access.SubjectType != base.SubjectTypeProject || access.SubjectProject == nil {
			continue
		}
		projResp := &SSHKeyProjectAccessResp{
			ID:      access.SubjectID,
			Name:    access.SubjectProject.Name,
			Allowed: access.Actions.Read || access.Actions.Write || access.Actions.Delete,
		}
		resp = append(resp, projResp)
		mapProjectResp[access.SubjectID] = projResp
	}
	for _, access := range accesses {
		if access.SubjectType != base.SubjectTypeApp || access.SubjectApp == nil {
			continue
		}
		projResp := mapProjectResp[access.SubjectApp.ProjectID]
		if projResp != nil {
			projResp.AppAccesses = append(projResp.AppAccesses, &SSHKeyAppAccessResp{
				ID:      access.SubjectID,
				Name:    access.SubjectApp.Name,
				Allowed: access.Actions.Read || access.Actions.Write || access.Actions.Delete,
			})
		}
	}
	return resp, nil
}
