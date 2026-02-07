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
	CronType     base.CronJobType          `json:"cronType"`
	CronExpr     string                    `json:"cronExpr"`
	App          ObjectID                  `json:"app,omitzero"`
	InitialTime  time.Time                 `json:"initialTime"`
	Priority     base.TaskPriority         `json:"priority"`
	MaxRetry     int                       `json:"maxRetry"`
	RetryDelay   timeutil.Duration         `json:"retryDelay"`
	Timeout      timeutil.Duration         `json:"timeout"`
	Command      *CronJobContainerCommand  `json:"command"`
	Notification *DefaultResultNtfnSetting `json:"notification,omitempty"`
}

type CronJobContainerCommand struct {
	Command    string `json:"command"`
	WorkingDir string `json:"workingDir,omitempty"`
}

func (s *CronJob) GetType() base.SettingType {
	return base.SettingTypeCronJob
}

func (s *CronJob) GetRefSettingIDs() []string {
	res := make([]string, 0, 5) //nolint
	res = append(res, s.Notification.GetRefSettingIDs()...)
	return res
}

func (s *CronJob) ParseCronExpr() (cron.Schedule, error) {
	sched, err := parser.Parse(s.CronExpr)
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
