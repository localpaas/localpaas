package discorddto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

const (
	maskedWebhook = "****************"
)

type GetDiscordReq struct {
	ID string `json:"-"`
}

func NewGetDiscordReq() *GetDiscordReq {
	return &GetDiscordReq{}
}

func (req *GetDiscordReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetDiscordResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *DiscordResp      `json:"data"`
}

type DiscordResp struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Webhook   string `json:"webhook"`
	Encrypted bool   `json:"encrypted,omitempty"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func TransformDiscord(setting *entity.Setting, decrypt bool) (resp *DiscordResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	config, err := setting.ParseDiscord(decrypt)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.IsEncrypted()
	if resp.Encrypted {
		resp.Webhook = maskedWebhook
	}
	return resp, nil
}
