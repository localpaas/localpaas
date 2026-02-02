package webhookuc

import (
	"errors"

	"github.com/gitsight/go-vcsurl"
	"github.com/go-playground/webhooks/v6/bitbucket"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/githelper"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

func (uc *WebhookUC) processBitbucketWebhook(
	req *webhookdto.HandleGitWebhookReq,
	secret string,
	apps []*entity.App,
	appsToRedeploy *[]*entity.App,
) (bool, error) {
	hook, err := bitbucket.New(bitbucket.Options.UUID(secret))
	if err != nil {
		return false, nil //nolint
	}
	payload, err := hook.Parse(req.Request, bitbucket.RepoPushEvent)
	if err != nil {
		if errors.Is(err, bitbucket.ErrEventNotFound) { // ok event wasn't one of the ones asked to be parsed
			return true, nil
		}
		return false, nil //nolint
	}

	switch payload.(type) { //nolint
	case bitbucket.RepoPushPayload:
		push, _ := payload.(bitbucket.RepoPushPayload) //nolint
		repoRef := string(githelper.NormalizeRepoRef(push.Push.Changes[0].New.Name))
		repoURL, err := vcsurl.Parse(push.Repository.Links.HTML.Href)
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
