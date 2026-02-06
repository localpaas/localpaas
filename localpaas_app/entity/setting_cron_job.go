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
	CronType     base.CronJobType         `json:"cronType"`
	CronExpr     string                   `json:"cronExpr"`
	App          ObjectID                 `json:"app,omitzero"`
	InitialTime  time.Time                `json:"initialTime"`
	Priority     base.TaskPriority        `json:"priority"`
	MaxRetry     int                      `json:"maxRetry"`
	RetryDelay   timeutil.Duration        `json:"retryDelay"`
	Timeout      timeutil.Duration        `json:"timeout"`
	Command      *CronJobContainerCommand `json:"command"`
	Notification *CronJobNtfnSettings     `json:"notification,omitempty"`
}

type CronJobContainerCommand struct {
	Command    string `json:"command"`
	WorkingDir string `json:"workingDir,omitempty"`
}

type CronJobNtfnSettings struct {
	Success *CronJobTargetNtfnSettings `json:"success,omitempty"`
	Failure *CronJobTargetNtfnSettings `json:"failure,omitempty"`
}

func (s *CronJobNtfnSettings) HasViaEmailNtfnSettings() bool {
	return (s.Success != nil && s.Success.ViaEmail != nil) || (s.Failure != nil && s.Failure.ViaEmail != nil)
}

func (s *CronJobNtfnSettings) HasViaSlackNtfnSettings() bool {
	return (s.Success != nil && s.Success.ViaSlack != nil) || (s.Failure != nil && s.Failure.ViaSlack != nil)
}

func (s *CronJobNtfnSettings) HasViaDiscordNtfnSettings() bool {
	return (s.Success != nil && s.Success.ViaDiscord != nil) || (s.Failure != nil && s.Failure.ViaDiscord != nil)
}

type CronJobTargetNtfnSettings struct {
	ViaEmail   *EmailNtfnSetting   `json:"viaEmail,omitempty"`
	ViaSlack   *SlackNtfnSetting   `json:"viaSlack,omitempty"`
	ViaDiscord *DiscordNtfnSetting `json:"viaDiscord,omitempty"`
}

func (s *CronJob) GetType() base.SettingType {
	return base.SettingTypeCronJob
}

func (s *CronJob) GetRefSettingIDs() []string {
	return nil
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
