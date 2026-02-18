package healthcheckdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type GetHealthcheckReq struct {
	settings.GetSettingReq
}

func NewGetHealthcheckReq() *GetHealthcheckReq {
	return &GetHealthcheckReq{}
}

func (req *GetHealthcheckReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetHealthcheckResp struct {
	Meta *basedto.Meta    `json:"meta"`
	Data *HealthcheckResp `json:"data"`
}

type HealthcheckResp struct {
	*settings.BaseSettingResp
	HealthcheckType base.HealthcheckType         `json:"healthcheckType"`
	Interval        timeutil.Duration            `json:"interval"`
	MaxRetry        int                          `json:"maxRetry"`
	RetryDelay      timeutil.Duration            `json:"retryDelay"`
	Timeout         timeutil.Duration            `json:"timeout"`
	SaveResultTasks bool                         `json:"saveResultTasks"`
	REST            *HealthcheckRESTReq          `json:"rest"`
	GRPC            *HealthcheckGRPCReq          `json:"grpc"`
	Notification    *HealthcheckNotificationResp `json:"notification"`
}

type HealthcheckRESTResp struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	ContentType string `json:"contentType"`
	Body        string `json:"body"`
	ReturnCode  int    `json:"returnCode"`
	ReturnText  string `json:"returnText"`
	ReturnJSON  string `json:"returnJSON"`
}

type HealthcheckGRPCResp struct {
	Version      base.HealthcheckGRPCVersion `json:"version"`
	Addr         string                      `json:"addr"`
	Service      string                      `json:"service"`
	ReturnStatus base.HealthcheckGRPCStatus  `json:"returnStatus"`
}

type HealthcheckNotificationResp struct {
	Success *settings.BaseSettingResp `json:"success"`
	Failure *settings.BaseSettingResp `json:"failure"`
}

func TransformHealthcheck(
	setting *entity.Setting,
	refObjects *entity.RefObjects,
) (resp *HealthcheckResp, err error) {
	config := setting.MustAsHealthcheck()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if resp.Notification.Success != nil {
		itemResp, _ := settings.TransformSettingBase(refObjects.RefSettings[resp.Notification.Success.ID])
		resp.Notification.Success = itemResp
	}
	if resp.Notification.Failure != nil {
		itemResp, _ := settings.TransformSettingBase(refObjects.RefSettings[resp.Notification.Failure.ID])
		resp.Notification.Failure = itemResp
	}

	return resp, nil
}
