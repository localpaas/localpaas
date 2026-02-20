package taskhealthcheck

import (
	"context"
	"fmt"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity/cacheentity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/strutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/service/notificationservice"
)

const (
	eventSuccess = "success"
	eventFailure = "failure"
)

//nolint:unparam
func (e *Executor) sendNotification(
	ctx context.Context,
	db database.IDB,
	data *taskData,
) error {
	if data.Healthcheck.Notification == nil {
		return nil
	}
	notifSettingID := gofn.If(data.Task.IsDone(), data.Healthcheck.Notification.Success.ID,
		data.Healthcheck.Notification.Failure.ID)
	notifSetting := data.RefObjects.RefSettings[notifSettingID]
	if notifSetting == nil {
		return nil
	}
	notification := notifSetting.MustAsNotification()
	minSendingInterval := notification.MinSendInterval.ToDuration()

	currEvent := gofn.If(data.Task.IsDone(), eventSuccess, eventFailure)
	lastNotifEvent := data.NotifEventMap[data.HealthcheckSetting.ID]
	timeNow := timeutil.NowUTC()
	shouldSkipNotifEmail := false
	shouldSkipNotifSlack := false
	shouldSkipNotifDiscord := false

	if minSendingInterval > 0 && lastNotifEvent != nil && lastNotifEvent.Event == currEvent {
		if timeNow.Sub(lastNotifEvent.Ts) < minSendingInterval {
			shouldSkipNotifEmail = lastNotifEvent.EmailSent
			shouldSkipNotifSlack = lastNotifEvent.SlackSent
			shouldSkipNotifDiscord = lastNotifEvent.DiscordSent
		}
	}

	var execFuncs []func(ctx context.Context) error
	emailSent := false
	slackSent := false
	discordSent := false

	if !shouldSkipNotifEmail && notification.HasNotificationViaEmail() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			err := e.sendNotificationViaEmail(ctx, db, data)
			emailSent = err == nil
			return err
		})
	}
	if !shouldSkipNotifSlack && notification.HasNotificationViaSlack() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			err := e.sendNotificationViaSlack(ctx, db, data)
			slackSent = err == nil
			return err
		})
	}
	if !shouldSkipNotifDiscord && notification.HasNotificationViaDiscord() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			err := e.sendNotificationViaDiscord(ctx, db, data)
			discordSent = err == nil
			return err
		})
	}
	if len(execFuncs) == 0 {
		return nil
	}

	e.buildNotificationMsgData(data)

	_ = gofn.ExecTasksEx(ctx, 20, false, execFuncs...) //nolint

	// Update notification events in redis
	if minSendingInterval > 0 {
		_ = e.notifEventRepo.Set(ctx, data.HealthcheckSetting.ID, &cacheentity.HealthcheckNotifEvent{
			Event:       currEvent,
			Ts:          timeNow,
			EmailSent:   emailSent,
			SlackSent:   slackSent,
			DiscordSent: discordSent,
		}, minSendingInterval)
	}

	return nil
}

func (e *Executor) buildNotificationMsgData(
	data *taskData,
) {
	msgData := &notificationservice.BaseMsgDataHealthcheckNotification{
		Succeeded:       data.Task.IsDone(),
		HealthcheckName: data.HealthcheckSetting.Name,
		HealthcheckType: data.Healthcheck.HealthcheckType,
		StartedAt:       data.Task.StartedAt,
		Duration:        data.Task.GetDuration(),
		Retries:         data.Task.Config.Retry,
		DashboardLink:   config.Current.DashboardHealthcheckDetailsURL(data.HealthcheckSetting.ID, data.Task.ID),
	}
	if data.Project != nil {
		msgData.ProjectName = data.Project.Name
	}
	if data.App != nil {
		msgData.AppName = data.App.Name
	}

	output, _ := data.Task.OutputAsHealthcheck()
	if output.REST != nil && data.Healthcheck.REST != nil {
		input := data.Healthcheck.REST
		maxLen := 100
		pad := "..."
		if output.REST.ReturnCode != 0 {
			msgData.Expect = fmt.Sprintf("Status code = %v", input.ReturnCode)
			msgData.Actual = fmt.Sprintf("Status code = %v", output.REST.ReturnCode)
		}
		if output.REST.ReturnText != "" {
			msgData.Expect = fmt.Sprintf("Text = %v", strutil.CutShort(input.ReturnText, maxLen, pad))
			msgData.Actual = fmt.Sprintf("Text = %v", strutil.CutShort(output.REST.ReturnText, maxLen, pad))
		}
		if output.REST.ReturnJSON != "" {
			msgData.Expect = fmt.Sprintf("JSON = %v", strutil.CutShort(input.ReturnJSON, maxLen, pad))
			msgData.Actual = fmt.Sprintf("JSON = %v", strutil.CutShort(output.REST.ReturnJSON, maxLen, pad))
		}
	}
	if output.GRPC != nil && data.Healthcheck.GRPC != nil {
		if output.GRPC.ReturnStatus != 0 {
			msgData.Expect = fmt.Sprintf("Status = %v", data.Healthcheck.GRPC.ReturnStatus)
			msgData.Actual = fmt.Sprintf("Status = %v", output.GRPC.ReturnStatus)
		}
	}

	data.NotifMsgData = msgData
}

func (e *Executor) sendNotificationViaEmail(
	ctx context.Context,
	db database.IDB,
	data *taskData,
) error {
	settingID := gofn.If(data.Task.IsDone(), data.Healthcheck.Notification.Success.ID,
		data.Healthcheck.Notification.Failure.ID)
	settings := data.RefObjects.RefSettings[settingID].MustAsNotification()
	if settings == nil || settings.ViaEmail == nil {
		return nil
	}

	emailSetting := data.RefObjects.RefSettings[settings.ViaEmail.Sender.ID]
	if emailSetting == nil {
		return apperrors.NewMissing("Sender email account")
	}
	emailAcc := emailSetting.MustAsEmail() //nolint
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
	subject += gofn.If(data.Task.IsDone(), " Healthcheck succeeded", " Healthcheck failed")

	err = e.notificationService.EmailSendHealthcheckNotification(ctx, db,
		&notificationservice.EmailMsgDataHealthcheckNotification{
			BaseMsgDataHealthcheckNotification: data.NotifMsgData,
			Email:                              emailAcc,
			Recipients:                         userEmails,
			Subject:                            subject,
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
	settingID := gofn.If(data.Task.IsDone(), data.Healthcheck.Notification.Success.ID,
		data.Healthcheck.Notification.Failure.ID)
	settings := data.RefObjects.RefSettings[settingID].MustAsNotification()
	if settings == nil || settings.ViaSlack == nil {
		return nil
	}

	imSetting := data.RefObjects.RefSettings[settings.ViaSlack.Webhook.ID]
	if imSetting == nil {
		return apperrors.NewMissing("Slack webhook")
	}
	imService := imSetting.MustAsIMService() //nolint
	if imService == nil || imService.Slack == nil {
		return apperrors.NewMissing("Slack webhook")
	}

	err := e.notificationService.SlackSendHealthcheckNotification(ctx, db,
		&notificationservice.SlackMsgDataHealthcheckNotification{
			BaseMsgDataHealthcheckNotification: data.NotifMsgData,
			Setting:                            imService.Slack,
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
	settingID := gofn.If(data.Task.IsDone(), data.Healthcheck.Notification.Success.ID,
		data.Healthcheck.Notification.Failure.ID)
	settings := data.RefObjects.RefSettings[settingID].MustAsNotification()
	if settings == nil || settings.ViaDiscord == nil {
		return nil
	}

	imSetting := data.RefObjects.RefSettings[settings.ViaDiscord.Webhook.ID]
	if imSetting == nil {
		return apperrors.NewMissing("Discord webhook")
	}
	imService := imSetting.MustAsIMService() //nolint
	if imService == nil || imService.Discord == nil {
		return apperrors.NewMissing("Discord webhook")
	}

	err := e.notificationService.DiscordSendHealthcheckNotification(ctx, db,
		&notificationservice.DiscordMsgDataHealthcheckNotification{
			BaseMsgDataHealthcheckNotification: data.NotifMsgData,
			Setting:                            imService.Discord,
		})
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
