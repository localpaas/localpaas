package cacheentity

import (
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/services/git/github"
)

type GithubAppManifest struct {
	Manifest    *github.AppManifest `json:"manifest"`
	State       string              `json:"state"`
	Reprovision bool                `json:"reprovision"`
	GithubApp   *entity.Setting     `json:"githubApp"`
}
