package slackdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

const (
	maskedWebhook = "****************"
)

type GetSlackReq struct {
	ID string `json:"-"`
}

func NewGetSlackReq() *GetSlackReq {
	return &GetSlackReq{}
}

func (req *GetSlackReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetSlackResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *SlackResp        `json:"data"`
}

type SlackResp struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Status    base.SettingStatus `json:"status"`
	Webhook   string             `json:"webhook"`
	Encrypted bool               `json:"encrypted,omitempty"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func TransformSlack(setting *entity.Setting) (resp *SlackResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	config := setting.MustAsSlack()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.IsEncrypted()
	if resp.Encrypted {
		resp.Webhook = maskedWebhook
	}
	return resp, nil
}
