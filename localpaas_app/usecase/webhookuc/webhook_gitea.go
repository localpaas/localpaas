package webhookuc

import (
	"errors"

	"github.com/go-playground/webhooks/v6/gitea"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

func (uc *WebhookUC) processGiteaWebhook(
	req *webhookdto.HandleGitWebhookReq,
	data *eventData,
) error {
	hook, err := gitea.New()
	if err != nil {
		return nil //nolint
	}
	payload, err := hook.Parse(req.Request, gitea.PushEvent)
	if err != nil {
		if errors.Is(err, gitea.ErrEventNotFound) { // ok event wasn't one of the ones asked to be parsed
			return nil
		}
		return apperrors.Wrap(err)
	}

	switch payload.(type) { //nolint
	case gitea.PushPayload:
		push, _ := payload.(gitea.PushPayload) //nolint
		data.Push = &pushEventData{
			RepoRef: push.Ref,
			RepoURL: push.Repo.HTMLURL,
		}
	}
	return nil
}
