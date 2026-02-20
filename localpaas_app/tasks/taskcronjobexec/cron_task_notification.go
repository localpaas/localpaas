package taskcronjobexec

import (
	"context"
	"fmt"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/service/notificationservice"
)

func (e *Executor) sendNotification(
	ctx context.Context,
	db database.IDB,
	data *taskData,
) error {
	if data.CronJob.Notification == nil {
		return nil
	}

	notifSettingID := gofn.If(data.Task.IsDone(), data.CronJob.Notification.Success.ID,
		data.CronJob.Notification.Failure.ID)
	notifSetting := data.RefObjects.RefSettings[notifSettingID]
	if notifSetting == nil {
		return nil
	}
	notification := notifSetting.MustAsNotification()

	var execFuncs []func(ctx context.Context) error

	if notification.HasNotificationViaEmail() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.sendNotificationViaEmail(ctx, db, data)
		})
	}
	if notification.HasNotificationViaSlack() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.sendNotificationViaSlack(ctx, db, data)
		})
	}
	if notification.HasNotificationViaDiscord() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.sendNotificationViaDiscord(ctx, db, data)
		})
	}
	if len(execFuncs) == 0 {
		return nil
	}

	e.buildNotificationMsgData(data)

	err := gofn.ExecTasks(ctx, 0, execFuncs...)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) buildNotificationMsgData(
	data *taskData,
) {
	msgData := &notificationservice.BaseMsgDataCronTaskNotification{
		Succeeded:     data.Task.IsDone(),
		CronJobName:   data.CronJobSetting.Name,
		CronJobExpr:   data.CronJob.CronExpr,
		CreatedAt:     data.CronJob.InitialTime,
		StartedAt:     data.Task.StartedAt,
		Duration:      data.Task.GetDuration(),
		Retries:       data.Task.Config.Retry,
		DashboardLink: config.Current.DashboardCronTaskDetailsURL(data.CronJobSetting.ID, data.Task.ID),
	}
	if data.Project != nil {
		msgData.ProjectName = data.Project.Name
	}
	if data.App != nil {
		msgData.AppName = data.App.Name
	}
	data.NotifMsgData = msgData
}

func (e *Executor) sendNotificationViaEmail(
	ctx context.Context,
	db database.IDB,
	data *taskData,
) error {
	settingID := gofn.If(data.Task.IsDone(), data.CronJob.Notification.Success.ID,
		data.CronJob.Notification.Failure.ID)
	settings := data.RefObjects.RefSettings[settingID].MustAsNotification()
	if settings == nil || settings.ViaEmail == nil {
		return nil
	}

	emailSetting := data.RefObjects.RefSettings[settings.ViaEmail.Sender.ID]
	if emailSetting == nil {
		return apperrors.NewMissing("Sender email account")
	}
	emailAcc := emailSetting.MustAsEmail()
	if emailAcc == nil {
		return apperrors.NewMissing("Sender email account")
	}

	userMap, err := e.userService.LoadProjectUsers(ctx, db, data.Project, settings.ViaEmail.ToProjectMembers,
		settings.ViaEmail.ToProjectOwners, settings.ViaEmail.ToAllAdmins)
	if err != nil {
		return apperrors.Wrap(err)
	}

	userEmails := make([]string, 0, len(userMap))
	for _, user := range userMap {
		userEmails = append(userEmails, user.Email)
	}
	if len(settings.ViaEmail.ToAddresses) > 0 {
		userEmails = gofn.ToSet(append(userEmails, settings.ViaEmail.ToAddresses...))
	}
	if len(userEmails) == 0 {
		return nil
	}

	subject := fmt.Sprintf("[%s/%s]", data.Project.Name, data.App.Name)
	if data.Task.IsDone() {
		subject += " scheduled task succeeded"
	} else {
		subject += " scheduled task failed"
	}

	err = e.notificationService.EmailSendCronTaskNotification(ctx, db,
		&notificationservice.EmailMsgDataCronTaskNotification{
			BaseMsgDataCronTaskNotification: data.NotifMsgData,
			Email:                           emailAcc,
			Recipients:                      userEmails,
			Subject:                         subject,
		})
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) sendNotificationViaSlack(
	ctx context.Context,
	db database.IDB,
	data *taskData,
) error {
	settingID := gofn.If(data.Task.IsDone(), data.CronJob.Notification.Success.ID,
		data.CronJob.Notification.Failure.ID)
	settings := data.RefObjects.RefSettings[settingID].MustAsNotification()
	if settings == nil || settings.ViaSlack == nil {
		return nil
	}

	imSetting := data.RefObjects.RefSettings[settings.ViaSlack.Webhook.ID]
	if imSetting == nil {
		return apperrors.NewMissing("Slack webhook")
	}
	imService := imSetting.MustAsIMService()
	if imService == nil || imService.Slack == nil {
		return apperrors.NewMissing("Slack webhook")
	}

	err := e.notificationService.SlackSendCronTaskNotification(ctx, db,
		&notificationservice.SlackMsgDataCronTaskNotification{
			BaseMsgDataCronTaskNotification: data.NotifMsgData,
			Setting:                         imService.Slack,
		})
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) sendNotificationViaDiscord(
	ctx context.Context,
	db database.IDB,
	data *taskData,
) error {
	settingID := gofn.If(data.Task.IsDone(), data.CronJob.Notification.Success.ID,
		data.CronJob.Notification.Failure.ID)
	settings := data.RefObjects.RefSettings[settingID].MustAsNotification()
	if settings == nil || settings.ViaDiscord == nil {
		return nil
	}

	imSetting := data.RefObjects.RefSettings[settings.ViaDiscord.Webhook.ID]
	if imSetting == nil {
		return apperrors.NewMissing("Discord webhook")
	}
	imService := imSetting.MustAsIMService()
	if imService == nil || imService.Discord == nil {
		return apperrors.NewMissing("Discord webhook")
	}

	err := e.notificationService.DiscordSendCronTaskNotification(ctx, db,
		&notificationservice.DiscordMsgDataCronTaskNotification{
			BaseMsgDataCronTaskNotification: data.NotifMsgData,
			Setting:                         imService.Discord,
		})
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
