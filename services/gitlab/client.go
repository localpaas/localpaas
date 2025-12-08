package gitlab

import (
	gogitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/httpclient"
)

type Client struct {
	token   string
	baseURL string

	client      *gogitlab.Client
	currentUser *gogitlab.User
}

func NewFromToken(token string, baseURL string) (*Client, error) {
	options := []gogitlab.ClientOptionFunc{
		gogitlab.WithHTTPClient(httpclient.DefaultClient),
	}
	if baseURL != "" {
		options = append(options, gogitlab.WithBaseURL(baseURL))
	}
	client, err := gogitlab.NewClient(token, options...)
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
		if gitSource != base.GitSourceGitlab && gitSource != base.GitSourceGitlabCustom {
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
