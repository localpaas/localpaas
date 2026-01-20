package discorddto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedWebhook = "****************"
)

type GetDiscordReq struct {
	settings.GetSettingReq
}

func NewGetDiscordReq() *GetDiscordReq {
	return &GetDiscordReq{}
}

func (req *GetDiscordReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetDiscordResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *DiscordResp      `json:"data"`
}

type DiscordResp struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Status    base.SettingStatus `json:"status"`
	Webhook   string             `json:"webhook"`
	Encrypted bool               `json:"encrypted,omitempty"`
	UpdateVer int                `json:"updateVer"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func (resp *DiscordResp) CopyWebhook(field entity.EncryptedField) error {
	resp.Webhook = field.String()
	return nil
}

func TransformDiscord(setting *entity.Setting) (resp *DiscordResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	config := setting.MustAsDiscord()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.Webhook.IsEncrypted()
	if resp.Encrypted {
		resp.Webhook = maskedWebhook
	}
	return resp, nil
}
