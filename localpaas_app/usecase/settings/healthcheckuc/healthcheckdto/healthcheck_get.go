package healthcheckdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/notification/notificationdto"
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
	HealthcheckType base.HealthcheckType                           `json:"healthcheckType"`
	Interval        timeutil.Duration                              `json:"interval"`
	MaxRetry        int                                            `json:"maxRetry"`
	RetryDelay      timeutil.Duration                              `json:"retryDelay"`
	Timeout         timeutil.Duration                              `json:"timeout"`
	REST            *HealthcheckRESTReq                            `json:"rest"`
	GRPC            *HealthcheckGRPCReq                            `json:"grpc"`
	Notification    *notificationdto.DefaultResultNotifSettingResp `json:"notification"`
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

type HealthcheckTransformInput struct {
	RefSettingMap map[string]*entity.Setting
}

func TransformHealthcheck(
	setting *entity.Setting,
	input *HealthcheckTransformInput,
) (resp *HealthcheckResp, err error) {
	config := setting.MustAsHealthcheck()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Notification, err = notificationdto.TransformDefaultResultNotifSetting(
		config.Notification, input.RefSettingMap)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}
