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
	CronType     base.CronJobType           `json:"cronType"`
	CronExpr     string                     `json:"cronExpr"`
	App          ObjectID                   `json:"app,omitzero"`
	InitialTime  time.Time                  `json:"initialTime"`
	Priority     base.TaskPriority          `json:"priority,omitempty"`
	MaxRetry     int                        `json:"maxRetry,omitempty"`
	RetryDelay   timeutil.Duration          `json:"retryDelay,omitempty"`
	Timeout      timeutil.Duration          `json:"timeout,omitempty"`
	Command      *CronJobContainerCommand   `json:"command,omitempty"`
	Notification *DefaultResultNotifSetting `json:"notification,omitempty"`
}

type CronJobContainerCommand struct {
	RunInShell string                    `json:"runInShell,omitempty"`
	Command    string                    `json:"command"`
	WorkingDir string                    `json:"workingDir,omitempty"`
	EnvVars    []*EnvVar                 `json:"envVars,omitempty"`
	ArgGroups  []*CronJobCommandArgGroup `json:"argGroups,omitempty"`
}

type CronJobCommandArgGroup struct {
	ExportEnv string               `json:"exportEnv"`
	Separator string               `json:"separator"`
	Args      []*CronJobCommandArg `json:"args,omitempty"`
}

type CronJobCommandArg struct {
	Use   bool   `json:"use"`
	Name  string `json:"name"`
	Value string `json:"value"`
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
