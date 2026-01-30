package entity

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

const (
	CurrentCronJobVersion = 1
)

var (
	parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
)

type CronJob struct {
	Cron        string            `json:"cron"`
	InitialTime time.Time         `json:"initialTime"`
	Priority    base.TaskPriority `json:"priority"`
	MaxRetry    int               `json:"maxRetry"`
	RetryDelay  timeutil.Duration `json:"retryDelay"`
	Timeout     timeutil.Duration `json:"timeout"`
	Command     string            `json:"command"`
}

func (s *CronJob) GetType() base.SettingType {
	return base.SettingTypeCronJob
}

func (s *CronJob) GetRefSettingIDs() []string {
	return nil
}

func (s *CronJob) ParseCron() (cron.Schedule, error) {
	sched, err := parser.Parse(s.Cron)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return sched, nil
}

func (s *Setting) AsCronJob() (*CronJob, error) {
	return parseSettingAs(s, func() *CronJob { return &CronJob{} })
}

func (s *Setting) MustAsCronJob() *CronJob {
	return gofn.Must(s.AsCronJob())
}
