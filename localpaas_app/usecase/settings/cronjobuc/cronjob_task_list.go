package cronjobuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/service/taskservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc/cronjobdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/taskuc/taskdto"
)

func (uc *CronJobUC) ListCronJobTask(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.ListCronJobTaskReq,
) (*cronjobdto.ListCronJobTaskResp, error) {
	req.Type = currentSettingType
	jobSetting, err := settings.GetSettingByID(ctx, uc.db, uc.settingRepo, &req.BaseSettingReq, req.JobID,
		false, false)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	listResp, err := uc.taskService.ListTask(ctx, uc.db, &taskservice.ListTaskReq{
		JobID:  []string{jobSetting.ID},
		Status: req.Status,
		Search: req.Search,
		Paging: req.Paging,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := taskdto.TransformTasks(listResp.Tasks, listResp.TaskInfoMap)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cronjobdto.ListCronJobTaskResp{
		Meta: &basedto.ListMeta{Page: listResp.PagingMeta},
		Data: resp,
	}, nil
}
