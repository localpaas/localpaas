package webhookuc

import (
	"errors"
	"net/http"

	"github.com/go-playground/webhooks/v6/gitea"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (uc *UC) parseGiteaWebhook(
	req *http.Request,
	secret string,
	data *repoEventData,
) error {
	hook, err := gitea.New(gitea.Options.Secret(secret))
	if err != nil {
		return apperrors.Wrap(err)
	}
	payload, err := hook.Parse(req, gitea.PushEvent)
	if err != nil {
		if errors.Is(err, gitea.ErrEventNotFound) { // ok event wasn't one of the ones asked to be parsed
			return nil
		}
		return apperrors.Wrap(err)
	}

	switch payload.(type) { //nolint
	case gitea.PushPayload:
		push, _ := payload.(gitea.PushPayload) //nolint
		data.Push = &repoPushEventData{
			RepoRef:  push.Ref,
			RepoURL:  push.Repo.HTMLURL,
			ChangeID: push.After,
		}
	}
	return nil
}
