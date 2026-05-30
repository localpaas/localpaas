package webhookuc

import (
	"errors"
	"net/http"

	"github.com/go-playground/webhooks/v6/github"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (uc *UC) parseGithubWebhook(
	req *http.Request,
	secret string,
	data *repoEventData,
) error {
	hook, err := github.New(github.Options.Secret(secret))
	if err != nil {
		return apperrors.Wrap(err)
	}
	payload, err := hook.Parse(req, github.PushEvent)
	if err != nil {
		if errors.Is(err, github.ErrEventNotFound) { // ok event wasn't one of the ones asked to be parsed
			return nil
		}
		return apperrors.Wrap(err)
	}

	switch payload.(type) { //nolint
	case github.PushPayload:
		push, _ := payload.(github.PushPayload) //nolint
		data.Push = &repoPushEventData{
			RepoRef:  push.Ref,
			RepoURL:  push.Repository.HTMLURL,
			ChangeID: push.After,
		}
	}
	return nil
}
