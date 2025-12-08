package gitlab

import (
	gogitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type Client struct {
	token string

	client      *gogitlab.Client
	currentUser *gogitlab.User
}

func NewFromToken(token string) (*Client, error) {
	client, err := gogitlab.NewClient(token)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return &Client{
		token:  token,
		client: client,
	}, nil
}

func NewFromSetting(setting *entity.Setting) (*Client, error) {
	switch setting.Type { //nolint:exhaustive
	case base.SettingTypeGitToken:
		gitToken, err := setting.AsGitToken()
		if base.GitSource(setting.Kind) != base.GitSourceGitlab {
			return nil, apperrors.New(ErrAccessProviderInvalid).
				WithMsgLog("git source '%s' is invalid", setting.Kind)
		}
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		token, err := gitToken.Token.GetPlain()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		return NewFromToken(token)

	default:
		return nil, apperrors.Wrap(ErrAccessProviderInvalid)
	}
}
