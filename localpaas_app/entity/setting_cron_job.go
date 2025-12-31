package entity

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentCronJobVersion = 1
)

var (
	parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
)

type CronJob struct {
	Cron         string            `json:"cron"`
	InitialTime  time.Time         `json:"initialTime"`
	Priority     base.TaskPriority `json:"priority"`
	MaxRetry     int               `json:"maxRetry"`
	RetryDelayMs int               `json:"retryDelayMs"`
	TimeoutMs    int               `json:"timeoutMs"`
	Command      string            `json:"command"`
}

func (j *CronJob) ParseCron() (cron.Schedule, error) {
	sched, err := parser.Parse(j.Cron)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return sched, nil
}

func (s *Setting) AsCronJob() (*CronJob, error) {
	return parseSettingAs(s, base.SettingTypeCronJob, func() *CronJob { return &CronJob{} })
}

func (s *Setting) MustAsCronJob() *CronJob {
	return gofn.Must(s.AsCronJob())
}
