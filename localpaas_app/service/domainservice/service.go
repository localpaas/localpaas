package domainservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type Service interface {
	VerifyProjectDomains(ctx context.Context, db database.IDB, projectID string, domains []string) error
	VerifyDomainsAvailable(ctx context.Context, db database.IDB, domains []string, ignoreAppIDs []string) error
}
