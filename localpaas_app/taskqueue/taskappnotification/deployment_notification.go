package taskappnotification

import (
	"context"
	"fmt"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/service/emailservice"
	"github.com/localpaas/localpaas/localpaas_app/service/imservice"
)

type deploymentNtfnTaskData struct {
	*taskData
}

func (e *Executor) notifyForDeployment(
	ctx context.Context,
	db database.Tx,
	taskData *taskData,
) error {
	data := &deploymentNtfnTaskData{taskData: taskData}
	ntfnSettings := data.NtfnSettings.Deployment
	var execFuncs []func(ctx context.Context) error

	if ntfnSettings.HasViaEmailNotificationSettings() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.notifyForDeploymentViaEmail(ctx, db, taskData)
		})
	}
	if ntfnSettings.HasViaSlackNotificationSettings() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.notifyForDeploymentViaSlack(ctx, db, taskData)
		})
	}
	if ntfnSettings.HasViaDiscordNotificationSettings() {
		execFuncs = append(execFuncs, func(ctx context.Context) error {
			return e.notifyForDeploymentViaDiscord(ctx, db, taskData)
		})
	}
	if len(execFuncs) == 0 {
		return nil
	}

	err := gofn.ExecTasks(ctx, 0, execFuncs...)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) notifyForDeploymentViaEmail(
	ctx context.Context,
	db database.Tx,
	data *taskData,
) error {
	success := data.Deployment.Status == base.DeploymentStatusDone
	settings := gofn.If(success, data.NtfnSettings.Deployment.Success, data.NtfnSettings.Deployment.Failure)
	if settings == nil || settings.ViaEmail == nil {
		return nil
	}

	senderEmail := data.RefSettingMap[settings.ViaEmail.Sender.ID]
	if senderEmail == nil {
		return apperrors.NewMissing("Sender email account")
	}

	userMap, err := e.userService.LoadProjectUsers(ctx, db, data.Project, settings.ViaEmail.ToProjectMembers,
		settings.ViaEmail.ToProjectOwners, settings.ViaEmail.ToAllAdmins)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if len(userMap) == 0 {
		return nil
	}

	userEmails := make([]string, 0, len(userMap))
	for _, user := range userMap {
		userEmails = append(userEmails, user.Email)
	}

	subject := fmt.Sprintf("[%s/%s]", data.Project.Name, data.App.Name)
	if success {
		subject += " deployment succeeded"
	} else {
		subject += " deployment failed"
	}

	deployment := data.Deployment
	emailData := &emailservice.EmailDataAppDeploymentNotification{
		Email:         senderEmail.MustAsEmail(),
		Recipients:    userEmails,
		Subject:       subject,
		ProjectName:   data.Project.Name,
		AppName:       data.App.Name,
		Succeeded:     success,
		Method:        deployment.Settings.ActiveMethod,
		Duration:      deployment.GetDuration(),
		DashboardLink: config.Current.DashboardDeploymentDetailsURL(deployment.ID),
	}

	switch deployment.Settings.ActiveMethod {
	case base.DeploymentMethodRepo:
		emailData.RepoURL = deployment.Settings.RepoSource.RepoURL
		emailData.RepoRef = deployment.Settings.RepoSource.RepoRef
	case base.DeploymentMethodImage:
		emailData.Image = deployment.Settings.ImageSource.Image
	case base.DeploymentMethodTarball:
		emailData.SourceArchive = "source tarball" //nolint TODO: update this
	}

	err = e.emailService.SendMailAppDeploymentNotification(ctx, db, emailData)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) notifyForDeploymentViaSlack(
	ctx context.Context,
	db database.Tx,
	data *taskData,
) error {
	success := data.Deployment.Status == base.DeploymentStatusDone
	settings := gofn.If(success, data.NtfnSettings.Deployment.Success, data.NtfnSettings.Deployment.Failure)
	if settings == nil || settings.ViaSlack == nil {
		return nil
	}

	imSetting := data.RefSettingMap[settings.ViaSlack.Webhook.ID]
	if imSetting == nil {
		return apperrors.NewMissing("Slack webhook")
	}
	imService := imSetting.MustAsIMService()
	if imService.Slack == nil {
		return apperrors.NewMissing("Slack webhook")
	}

	deployment := data.Deployment
	slackData := &imservice.SlackMsgDataAppDeploymentNotification{
		Slack:         imService.Slack,
		ProjectName:   data.Project.Name,
		AppName:       data.App.Name,
		Succeeded:     success,
		Method:        deployment.Settings.ActiveMethod,
		Duration:      deployment.GetDuration(),
		DashboardLink: config.Current.DashboardDeploymentDetailsURL(deployment.ID),
	}

	switch deployment.Settings.ActiveMethod {
	case base.DeploymentMethodRepo:
		slackData.RepoURL = deployment.Settings.RepoSource.RepoURL
		slackData.RepoRef = deployment.Settings.RepoSource.RepoRef
	case base.DeploymentMethodImage:
		slackData.Image = deployment.Settings.ImageSource.Image
	case base.DeploymentMethodTarball:
		slackData.SourceArchive = "source tarball" // TODO: update this
	}

	err := e.imService.SlackSendAppDeploymentNotification(ctx, db, slackData)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) notifyForDeploymentViaDiscord(
	ctx context.Context,
	db database.Tx,
	data *taskData,
) error {
	success := data.Deployment.Status == base.DeploymentStatusDone
	settings := gofn.If(success, data.NtfnSettings.Deployment.Success, data.NtfnSettings.Deployment.Failure)
	if settings == nil || settings.ViaSlack == nil {
		return nil
	}

	imSetting := data.RefSettingMap[settings.ViaDiscord.Webhook.ID]
	if imSetting == nil {
		return apperrors.NewMissing("Discord webhook")
	}
	imService := imSetting.MustAsIMService()
	if imService.Discord == nil {
		return apperrors.NewMissing("Discord webhook")
	}

	deployment := data.Deployment
	discordData := &imservice.DiscordMsgDataAppDeploymentNotification{
		Discord:       imService.Discord,
		ProjectName:   data.Project.Name,
		AppName:       data.App.Name,
		Succeeded:     success,
		Method:        deployment.Settings.ActiveMethod,
		Duration:      deployment.GetDuration(),
		DashboardLink: config.Current.DashboardDeploymentDetailsURL(deployment.ID),
	}

	switch deployment.Settings.ActiveMethod {
	case base.DeploymentMethodRepo:
		discordData.RepoURL = deployment.Settings.RepoSource.RepoURL
		discordData.RepoRef = deployment.Settings.RepoSource.RepoRef
	case base.DeploymentMethodImage:
		discordData.Image = deployment.Settings.ImageSource.Image
	case base.DeploymentMethodTarball:
		discordData.SourceArchive = "source tarball" // TODO: update this
	}

	err := e.imService.DiscordSendAppDeploymentNotification(ctx, db, discordData)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
