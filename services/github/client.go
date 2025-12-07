package github

import (
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v75/github"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

type Client struct {
	appID          int64
	installationID int64
	personalToken  string

	appsTransport    *ghinstallation.AppsTransport
	installTransport *ghinstallation.Transport

	client *github.Client
}

func (c *Client) isAppClient() bool {
	return c.appID > 0
}

func NewFromApp(appID, installationID int64, privateKey []byte) (*Client, error) {
	appTr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, appID, privateKey)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	client := &Client{
		appID:          appID,
		installationID: installationID,
		appsTransport:  appTr,
	}
	if installationID != 0 {
		client.installTransport = ghinstallation.NewFromAppsTransport(appTr, installationID)
		client.client = github.NewClient(&http.Client{Transport: client.installTransport})
	} else {
		client.client = github.NewClient(&http.Client{Transport: client.appsTransport})
	}

	return client, nil
}

func NewFromPersonalToken(personalToken string) (*Client, error) {
	client := &Client{
		personalToken: personalToken,
		client:        github.NewClient(&http.Client{Transport: NewPatTransport(http.DefaultTransport, personalToken)}),
	}
	return client, nil
}

func NewFromSetting(setting *entity.Setting) (*Client, error) {
	switch setting.Type { //nolint:exhaustive
	case base.SettingTypeGithubApp:
		githubApp, err := setting.AsGithubApp()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		privateKey, err := githubApp.PrivateKey.GetPlain()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		return NewFromApp(githubApp.AppID, githubApp.InstallationID, reflectutil.UnsafeStrToBytes(privateKey))

	case base.SettingTypeGitToken:
		gitToken, err := setting.AsGitToken()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		token, err := gitToken.Token.GetPlain()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		return NewFromPersonalToken(token)

	default:
		return nil, apperrors.Wrap(ErrGithubAccessProviderInvalid)
	}
}
