package webhookuc

import (
	"errors"

	"github.com/go-playground/webhooks/v6/bitbucket"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/githelper"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

func (uc *WebhookUC) processBitbucketWebhook(
	req *webhookdto.HandleGitWebhookReq,
	data *eventData,
) error {
	hook, err := bitbucket.New()
	if err != nil {
		return nil //nolint
	}
	payload, err := hook.Parse(req.Request, bitbucket.RepoPushEvent)
	if err != nil {
		if errors.Is(err, bitbucket.ErrEventNotFound) { // ok event wasn't one of the ones asked to be parsed
			return nil
		}
		return apperrors.Wrap(err)
	}

	switch payload.(type) { //nolint
	case bitbucket.RepoPushPayload:
		push, _ := payload.(bitbucket.RepoPushPayload) //nolint
		data.Push = &pushEventData{
			RepoRef: string(githelper.NormalizeRepoRef(push.Push.Changes[0].New.Name)),
			RepoURL: push.Repository.Links.HTML.Href,
		}
	}
	return nil
}
