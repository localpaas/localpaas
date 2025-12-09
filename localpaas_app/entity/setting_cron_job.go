package entity

import (
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type CronJob struct {
	Cron        string    `json:"cron"`
	InitialTime time.Time `json:"initialTime"`
	Command     string    `json:"command"`
}

func (s *Setting) AsCronJob() (*CronJob, error) {
	return parseSettingAs(s, base.SettingTypeCronJob, func() *CronJob { return &CronJob{} })
}

func (s *Setting) MustAsCronJob() *CronJob {
	return gofn.Must(s.AsCronJob())
}
