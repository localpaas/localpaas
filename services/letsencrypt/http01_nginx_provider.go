package letsencrypt

import (
	"os"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	// 0755 grants read/write/execute for owner, read/execute for group/others
	// 0644 grants read/write for owner, read-only for group/others
	http01DirFileMode = 0644
)

type Http01NginxProvider struct {
	http01ChallengeRoot string
}

func NewHttp01NginxProvider(http01ChallengeRoot string) *Http01NginxProvider {
	return &Http01NginxProvider{http01ChallengeRoot: http01ChallengeRoot}
}

// Present starts serving the token at `ChallengePath(token)` via nginx http-01 root
func (s *Http01NginxProvider) Present(domain, token, keyAuth string) error {
	err := os.MkdirAll(s.http01ChallengeRoot, http01DirFileMode)
	if err != nil {
		return apperrors.New(err).WithMsgLog("unable to create http-01-challenge directory")
	}

	file, err := os.Create(s.http01ChallengeRoot + "/" + token)
	if err != nil {
		return apperrors.New(err).WithMsgLog("unable to create token file")
	}

	_, err = file.Write([]byte(keyAuth))
	if err != nil {
		return apperrors.New(err).WithMsgLog("unable to write data to token file")
	}

	return nil
}

func (s *Http01NginxProvider) CleanUp(domain, token, _ string) error {
	err := os.Remove(s.http01ChallengeRoot + "/" + token)
	if err != nil {
		return apperrors.New(err).WithMsgLog("unable to clean token file after the challenge")
	}
	return nil
}
