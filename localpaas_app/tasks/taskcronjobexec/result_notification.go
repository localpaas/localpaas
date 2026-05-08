package taskcronjobexec

import (
	"context"
	"fmt"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/service/notificationservice"
)

func (e *Executor) sendNotification(
	ctx context.Context,
	db database.IDB,
	data *taskData,
) (err error) {
	cronJob := data.CronJob.MustAsCronJob()
	notifConfig := cronJob.Notification
	if notifConfig == nil {
		return nil
	}

	var scope *base.SettingScope
	switch {
	case data.App != nil:
		scope = data.App.GetSettingScope()
	case data.Project != nil:
		scope = data.Project.GetSettingScope()
	default:
		scope = base.NewSettingScopeGlobal()
	}

	isSucceeded := data.Task.IsDone()
	notification, err := e.notificationService.GetNotificationForEvent(ctx, db,
		scope, notifConfig, isSucceeded, data.RefObjects)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if notification == nil {
		return nil
	}

	e.buildNotificationMsgData(data)
	_, err = e.notificationService.NotifyForTaskResult(ctx, db, &notificationservice.TaskResultNotificationReq{
		ActionSucceeded: isSucceeded,
		ScopeProject:    data.Project,
		ScopeApp:        data.App,
		RefObjects:      data.RefObjects,
		Notification:    notification,
		TemplateName:    notificationservice.TemplateCronTaskNotification,
		TemplateData:    data.NotifMsgData,
	})
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (e *Executor) buildNotificationMsgData(
	data *taskData,
) {
	cronJob := data.CronJob.MustAsCronJob()
	isSucceeded := data.Task.IsDone()
	msgData := &notificationservice.TemplateDataCronTask{
		BaseTemplateData: notificationservice.BaseTemplateData{
			Title: e.notificationService.BuildTitlePrefix(data.Project, data.App, nil) +
				gofn.If(isSucceeded, " Scheduled task succeeded", " Scheduled task failed"),
		},
		Succeeded:   isSucceeded,
		CronJobName: data.CronJob.Name,
		CreatedAt:   cronJob.Schedule.InitialTime,
		StartedAt:   data.Task.StartedAt,
		Duration:    data.Task.GetDuration(),
		Retries:     data.Task.Config.Retry,
	}
	if cronJob.Schedule.Interval > 0 {
		msgData.Schedule = fmt.Sprintf("every %v", cronJob.Schedule.Interval.String())
	} else {
		msgData.Schedule = fmt.Sprintf("cron expression %v", cronJob.Schedule.CronExpr)
	}
	if data.Project != nil {
		msgData.ProjectName = data.Project.Name
	}
	if data.App != nil {
		msgData.AppName = data.App.Name
	}
	switch {
	case data.App != nil:
		msgData.DashboardLink = config.Current.DashboardAppCronTaskDetailsURL(data.App.ID, data.App.ProjectID,
			data.CronJob.ID, data.Task.ID)
	case data.Project != nil:
		msgData.DashboardLink = config.Current.DashboardProjectCronTaskDetailsURL(data.Project.ID,
			data.CronJob.ID, data.Task.ID)
	default:
		msgData.DashboardLink = config.Current.DashboardGlobalCronTaskDetailsURL(
			data.CronJob.ID, data.Task.ID)
	}
	data.NotifMsgData = msgData
}
