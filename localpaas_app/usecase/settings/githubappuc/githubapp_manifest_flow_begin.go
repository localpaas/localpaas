package githubappuc

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/entity/cacheentity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc/githubappdto"
	"github.com/localpaas/localpaas/services/git/github"
)

const (
	appManifestStateLen = 24
	appManifestCacheExp = 10 * time.Minute
)

func (uc *GithubAppUC) BeginGithubAppManifestFlow(
	ctx context.Context,
	auth *basedto.Auth,
	req *githubappdto.BeginGithubAppManifestFlowReq,
) (*githubappdto.BeginGithubAppManifestFlowResp, error) {
	cfg := config.Current
	isLocalEnv := cfg.IsDevEnv() && cfg.Platform == config.PlatformLocal
	timeNow := timeutil.NowUTC()

	appSetting := &entity.Setting{
		ID:              gofn.Must(ulid.NewStringULID()),
		ObjectID:        req.ObjectID,
		Type:            base.SettingTypeGithubApp,
		Kind:            string(base.SettingTypeGithubApp),
		Status:          base.SettingStatusActive,
		Name:            gofn.Coalesce(req.Name, "my localpaas app"),
		AvailInProjects: req.AvailInProjects,
		Default:         req.Default,
		Version:         0,
		CreatedAt:       timeNow,
		UpdatedAt:       timeNow,
	}
	githubApp := &entity.GithubApp{
		Organization: req.Org,
		SSOEnabled:   true,
	}
	if isLocalEnv {
		githubApp.WebhookSecret = "abc123"
		githubApp.WebhookURL = "https://smee.io/RBNiNjxieUIWZ6Ej"
	} else {
		githubApp.WebhookSecret = gofn.RandTokenAsHex(webhookSecretLen)
		githubApp.WebhookURL = cfg.RepoWebhookURL(base.WebhookKindGithub, githubApp.WebhookURL)
	}
	appSetting.MustSetData(githubApp)

	state := gofn.RandTokenAsHex(appManifestStateLen)
	manifest := &github.AppManifest{
		Name:         appSetting.Name,
		URL:          cfg.BaseURL,
		CallbackURLs: []string{cfg.SsoCallbackURL(appSetting.ID)},
		Hook: &github.AppManifestHook{
			URL:    githubApp.WebhookURL,
			Active: true,
		},
		Public: false,
		DefaultEvents: []string{
			"push",
			// "create",
		},
		DefaultPermissions: map[string]string{
			"contents": "read",
			// "repository_hooks": "write",
			// "organization_hooks": "write",
			// "administration": "read",
			// "issues": "read",
			// "repository_projects": "read",
			// "pull_requests": "read",
			// "organization_personal_access_tokens": "read",
		},
		RequestOAuthOnInstall: false,
		SetupOnUpdate:         true,
	}

	var beginFlowURL string
	switch req.Scope { //nolint:exhaustive
	case base.SettingScopeGlobal:
		beginFlowURL = cfg.GlobalGithubAppManifestFlowCreationURL(appSetting.ID, state)
		manifest.RedirectURL = cfg.GlobalGithubAppManifestFlowSetupURL(appSetting.ID)
		manifest.SetupURL = manifest.RedirectURL
	case base.SettingScopeProject:
		beginFlowURL = cfg.ProjectGithubAppManifestFlowCreationURL(req.ObjectID, appSetting.ID, state)
		manifest.RedirectURL = cfg.ProjectGithubAppManifestFlowSetupURL(req.ObjectID, appSetting.ID)
		manifest.SetupURL = manifest.RedirectURL
	default:
		return nil, apperrors.New(apperrors.ErrUnsupported)
	}

	manifestCache := &cacheentity.GithubAppManifest{
		Manifest:    manifest,
		State:       state,
		CreatingApp: appSetting,
	}

	err := uc.cacheAppManifestRepo.Set(ctx, appSetting.ID, manifestCache, appManifestCacheExp)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &githubappdto.BeginGithubAppManifestFlowResp{
		Data: &githubappdto.BeginGithubAppManifestFlowDataResp{
			RedirectURL: beginFlowURL,
		},
	}, nil
}
