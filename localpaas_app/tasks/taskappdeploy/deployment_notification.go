package taskappdeploy

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

func (e *Executor) notifyForDeployment(
	ctx context.Context,
	db database.IDB,
	data *taskData,
) error {
	if data.NotifSettings == nil || !data.NotifSettings.HasDeploymentNotifSetting() {
		return nil
	}
	notifSettings := data.NotifSettings.Deployment
	var execFuncs []func(ctx context.Context) error

	if notifSettings.HasViaEmailNotifSetting() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.notifyForDeploymentViaEmail(ctx, db, data)
		})
	}
	if notifSettings.HasViaSlackNotifSetting() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.notifyForDeploymentViaSlack(ctx, db, data)
		})
	}
	if notifSettings.HasViaDiscordNotifSetting() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.notifyForDeploymentViaDiscord(ctx, db, data)
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
	data *taskData,
) error {
	settings := gofn.If(data.Deployment.IsDone(), data.NotifSettings.Deployment.Success,
		data.NotifSettings.Deployment.Failure)
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
	data *taskData,
) error {
	settings := gofn.If(data.Deployment.IsDone(), data.NotifSettings.Deployment.Success,
		data.NotifSettings.Deployment.Failure)
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
	data *taskData,
) error {
	settings := gofn.If(data.Deployment.IsDone(), data.NotifSettings.Deployment.Success,
		data.NotifSettings.Deployment.Failure)
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
