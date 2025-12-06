package github

import (
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v75/github"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type App struct {
	appID          int64
	installationID int64

	appsTransport    *ghinstallation.AppsTransport
	installTransport *ghinstallation.Transport

	client *github.Client
}

func NewApp(appID, installationID int64, privateKey []byte) (*App, error) {
	appTr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, appID, privateKey)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	app := &App{
		appID:          appID,
		installationID: installationID,
		appsTransport:  appTr,
	}
	if installationID != 0 {
		app.installTransport = ghinstallation.NewFromAppsTransport(appTr, installationID)
		app.client = github.NewClient(&http.Client{Transport: app.installTransport})
	} else {
		app.client = github.NewClient(&http.Client{Transport: app.appsTransport})
	}

	return app, nil
}
