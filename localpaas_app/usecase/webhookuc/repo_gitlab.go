package webhookuc

import (
	"errors"
	"net/http"

	"github.com/go-playground/webhooks/v6/gitlab"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (uc *UC) parseGitlabWebhook(
	req *http.Request,
	secret string,
	data *repoEventData,
) error {
	hook, err := gitlab.New(gitlab.Options.Secret(secret))
	if err != nil {
		return apperrors.New(err)
	}
	payload, err := hook.Parse(req, gitlab.PushEvents)
	if err != nil {
		if errors.Is(err, gitlab.ErrEventNotFound) { // ok event wasn't one of the ones asked to be parsed
			return nil
		}
		return apperrors.New(err)
	}

	switch payload.(type) { //nolint
	case gitlab.PushEventPayload:
		push, _ := payload.(gitlab.PushEventPayload) //nolint
		data.Push = &repoPushEventData{
			RepoRef:  push.Ref,
			RepoURL:  push.Repository.GitHTTPURL,
			ChangeID: push.After,
		}
	}
	return nil
}
