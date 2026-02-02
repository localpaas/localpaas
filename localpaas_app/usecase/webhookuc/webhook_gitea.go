package webhookuc

import (
	"errors"

	"github.com/gitsight/go-vcsurl"
	"github.com/go-playground/webhooks/v6/gitea"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

func (uc *WebhookUC) processGiteaWebhook(
	req *webhookdto.HandleGitWebhookReq,
	secret string,
	apps []*entity.App,
	appsToRedeploy *[]*entity.App,
) (bool, error) {
	hook, err := gitea.New(gitea.Options.Secret(secret))
	if err != nil {
		return false, nil //nolint
	}
	payload, err := hook.Parse(req.Request, gitea.PushEvent)
	if err != nil {
		if errors.Is(err, gitea.ErrEventNotFound) { // ok event wasn't one of the ones asked to be parsed
			return true, nil
		}
		return false, nil //nolint
	}

	switch payload.(type) { //nolint
	case gitea.PushPayload:
		push, _ := payload.(gitea.PushPayload) //nolint
		repoRef := push.Ref
		repoURL, err := vcsurl.Parse(push.Repo.HTMLURL)
		if err != nil {
			return false, apperrors.Wrap(err)
		}
		for _, app := range apps {
			if flag, _ := uc.shouldRedeployApp(app, repoURL, repoRef); flag {
				*appsToRedeploy = append(*appsToRedeploy, app)
			}
		}
	}
	return true, nil
}
