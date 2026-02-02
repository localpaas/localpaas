package gitea

import (
	gogitea "code.gitea.io/sdk/gitea"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type Client struct {
	token   string
	baseURL string

	client *gogitea.Client
}

func NewFromToken(token string, baseURL string) (*Client, error) {
	client, err := gogitea.NewClient(baseURL, gogitea.SetToken(token))
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return &Client{
		token:   token,
		baseURL: baseURL,
		client:  client,
	}, nil
}

func NewFromSetting(setting *entity.Setting) (*Client, error) {
	switch setting.Type { //nolint:exhaustive
	case base.SettingTypeGitToken:
		gitToken, err := setting.AsGitToken()
		gitSource := base.GitSource(setting.Kind)
		if gitSource != base.GitSourceGitea {
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
		return NewFromToken(token, gitToken.BaseURL)

	default:
		return nil, apperrors.Wrap(ErrAccessProviderInvalid)
	}
}
