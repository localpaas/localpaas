package slackdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
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
	ID      string `json:"id"`
	Name    string `json:"name"`
	Webhook string `json:"webhook"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func TransformSlack(setting *entity.Setting) (resp *SlackResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	config, err := setting.ParseSlack()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}
