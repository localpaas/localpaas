package taskappdeploy

import (
	"context"
	"fmt"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/service/notificationservice"
)

func (e *Executor) notifyForDeployment(
	ctx context.Context,
	db database.IDB,
	data *taskData,
) error {
	if data.Deployment.Settings.Notification == nil {
		return nil
	}
	notifSettingID := gofn.If(data.Task.IsDone(), data.Deployment.Settings.Notification.Success.ID,
		data.Deployment.Settings.Notification.Failure.ID)
	notifSetting := data.RefObjects.RefSettings[notifSettingID]
	if notifSetting == nil {
		return nil
	}
	notification := notifSetting.MustAsNotification()

	var execFuncs []func(ctx context.Context) error

	if notification.HasNotificationViaEmail() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.notifyForDeploymentViaEmail(ctx, db, notification, data)
		})
	}
	if notification.HasNotificationViaSlack() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.notifyForDeploymentViaSlack(ctx, db, notification, data)
		})
	}
	if notification.HasNotificationViaDiscord() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.notifyForDeploymentViaDiscord(ctx, db, notification, data)
		})
	}
	if len(execFuncs) == 0 {
		return nil
	}

	e.buildDeploymentNotifMsgData(data)

	err := gofn.ExecTasks(ctx, 0, execFuncs...)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) buildDeploymentNotifMsgData(
	data *taskData,
) {
	deployment := data.Deployment

	msgData := &notificationservice.BaseMsgDataAppDeploymentNotification{
		ProjectName:   data.Project.Name,
		AppName:       data.App.Name,
		Succeeded:     deployment.IsDone(),
		Method:        deployment.Settings.ActiveMethod,
		StartedAt:     deployment.StartedAt,
		Duration:      deployment.GetDuration(),
		DashboardLink: config.Current.DashboardDeploymentDetailsURL(deployment.ID),
	}
	data.NotifMsgData = msgData

	switch deployment.Settings.ActiveMethod {
	case base.DeploymentMethodRepo:
		msgData.RepoURL = deployment.Settings.RepoSource.RepoURL
		msgData.RepoRef = deployment.Settings.RepoSource.RepoRef
		if deployment.Output != nil {
			msgData.CommitMsg = deployment.Output.CommitMessage
		}
	case base.DeploymentMethodImage:
		msgData.Image = deployment.Settings.ImageSource.Image
	}
}

func (e *Executor) notifyForDeploymentViaEmail(
	ctx context.Context,
	db database.IDB,
	notification *entity.Notification,
	data *taskData,
) error {
	if notification == nil || notification.ViaEmail == nil {
		return nil
	}

	emailSetting := data.RefObjects.RefSettings[notification.ViaEmail.Sender.ID]
	if emailSetting == nil {
		return apperrors.NewMissing("Sender email account")
	}
	emailAcc := emailSetting.MustAsEmail()
	if emailAcc == nil {
		return apperrors.NewMissing("Sender email account")
	}

	userMap, err := e.userService.LoadProjectUsers(ctx, db, data.Project, notification.ViaEmail.ToProjectMembers,
		notification.ViaEmail.ToProjectOwners, notification.ViaEmail.ToAllAdmins)
	if err != nil {
		return apperrors.Wrap(err)
	}

	userEmails := make([]string, 0, len(userMap))
	for _, user := range userMap {
		userEmails = append(userEmails, user.Email)
	}
	if len(notification.ViaEmail.ToAddresses) > 0 {
		userEmails = gofn.ToSet(append(userEmails, notification.ViaEmail.ToAddresses...))
	}
	if len(userEmails) == 0 {
		return nil
	}

	subject := fmt.Sprintf("[%s/%s]", data.Project.Name, data.App.Name)
	if data.Deployment.IsDone() {
		subject += " deployment succeeded"
	} else {
		subject += " deployment failed"
	}

	err = e.notificationService.EmailSendAppDeploymentNotification(ctx, db,
		&notificationservice.EmailMsgDataAppDeploymentNotification{
			BaseMsgDataAppDeploymentNotification: data.NotifMsgData,
			Email:                                emailAcc,
			Recipients:                           userEmails,
			Subject:                              subject,
		})
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) notifyForDeploymentViaSlack(
	ctx context.Context,
	db database.IDB,
	notification *entity.Notification,
	data *taskData,
) error {
	if notification == nil || notification.ViaSlack == nil {
		return nil
	}

	imSetting := data.RefObjects.RefSettings[notification.ViaSlack.Webhook.ID]
	if imSetting == nil {
		return apperrors.NewMissing("Slack webhook")
	}
	imService := imSetting.MustAsIMService()
	if imService == nil || imService.Slack == nil {
		return apperrors.NewMissing("Slack webhook")
	}

	err := e.notificationService.SlackSendAppDeploymentNotification(ctx, db,
		&notificationservice.SlackMsgDataAppDeploymentNotification{
			BaseMsgDataAppDeploymentNotification: data.NotifMsgData,
			Setting:                              imService.Slack,
		})
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) notifyForDeploymentViaDiscord(
	ctx context.Context,
	db database.IDB,
	notification *entity.Notification,
	data *taskData,
) error {
	if notification == nil || notification.ViaDiscord == nil {
		return nil
	}

	imSetting := data.RefObjects.RefSettings[notification.ViaDiscord.Webhook.ID]
	if imSetting == nil {
		return apperrors.NewMissing("Discord webhook")
	}
	imService := imSetting.MustAsIMService()
	if imService == nil || imService.Discord == nil {
		return apperrors.NewMissing("Discord webhook")
	}

	err := e.notificationService.DiscordSendAppDeploymentNotification(ctx, db,
		&notificationservice.DiscordMsgDataAppDeploymentNotification{
			BaseMsgDataAppDeploymentNotification: data.NotifMsgData,
			Setting:                              imService.Discord,
		})
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
