package sessionuc

import (
	"context"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gitea"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
)

func (uc *SessionUC) InitOAuthProvider(
	ctx context.Context,
	req *sessiondto.InitOAuthProviderReq,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, "", req.Name, true)
	if err != nil {
		return apperrors.Wrap(err)
	}

	oauth := setting.MustAsOAuth()
	clientSecret := oauth.ClientSecret.MustGetPlain()
	callbackURL := config.Current.SsoBaseCallbackURL() + "/" + req.Name

	var provider goth.Provider
	switch base.OAuthType(setting.Kind) {
	case base.OAuthTypeGithub, base.OAuthTypeGithubApp:
		provider = github.New(oauth.ClientID, clientSecret, callbackURL, oauth.Scopes...)
	case base.OAuthTypeGitlab:
		provider = gitlab.New(oauth.ClientID, clientSecret, callbackURL, oauth.Scopes...)
	case base.OAuthTypeGitea:
		provider = gitea.New(oauth.ClientID, clientSecret, callbackURL, oauth.Scopes...)
	case base.OAuthTypeGoogle:
		provider = google.New(oauth.ClientID, clientSecret, callbackURL, oauth.Scopes...)

	// Custom types
	case base.OAuthTypeGitlabCustom:
		provider = gitlab.NewCustomisedURL(oauth.ClientID, clientSecret, callbackURL,
			oauth.AuthURL, oauth.TokenURL, oauth.ProfileURL, oauth.Scopes...)
	}
	provider.SetName(req.Name)
	goth.UseProviders(provider)

	return nil
}
