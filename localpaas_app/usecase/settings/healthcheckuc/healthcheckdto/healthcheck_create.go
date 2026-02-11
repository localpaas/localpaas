package healthcheckdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/notification/notificationdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateHealthcheckReq struct {
	settings.CreateSettingReq
	*HealthcheckBaseReq
}

type HealthcheckBaseReq struct {
	Name            string                                        `json:"name"`
	HealthcheckType base.HealthcheckType                          `json:"healthcheckType"`
	Interval        timeutil.Duration                             `json:"interval"`
	MaxRetry        int                                           `json:"maxRetry"`
	RetryDelay      timeutil.Duration                             `json:"retryDelay"`
	Timeout         timeutil.Duration                             `json:"timeout"`
	REST            *HealthcheckRESTReq                           `json:"rest"`
	GRPC            *HealthcheckGRPCReq                           `json:"grpc"`
	Notification    *notificationdto.DefaultResultNotifSettingReq `json:"notification"`
}

func (req *HealthcheckBaseReq) ToEntity() *entity.Healthcheck {
	return &entity.Healthcheck{
		HealthcheckType: req.HealthcheckType,
		Interval:        req.Interval,
		MaxRetry:        req.MaxRetry,
		RetryDelay:      req.RetryDelay,
		Timeout:         req.Timeout,
		REST:            req.REST.ToEntity(),
		GRPC:            req.GRPC.ToEntity(),
		Notification:    req.Notification.ToEntity(),
	}
}

type HealthcheckRESTReq struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	ContentType string `json:"contentType"`
	Body        string `json:"body"`
	ReturnCode  int    `json:"returnCode"`
	ReturnText  string `json:"returnText"`
	ReturnJSON  string `json:"returnJSON"`
}

func (req *HealthcheckRESTReq) ToEntity() *entity.HealthcheckREST {
	if req == nil {
		return nil
	}
	return &entity.HealthcheckREST{
		URL:         req.URL,
		Method:      req.Method,
		ContentType: req.ContentType,
		Body:        req.Body,
		ReturnCode:  req.ReturnCode,
		ReturnText:  req.ReturnText,
		ReturnJSON:  req.ReturnJSON,
	}
}

type HealthcheckGRPCReq struct {
	Version      base.HealthcheckGRPCVersion `json:"version"`
	Addr         string                      `json:"addr"`
	Service      string                      `json:"service"`
	ReturnStatus base.HealthcheckGRPCStatus  `json:"returnStatus"`
}

func (req *HealthcheckGRPCReq) ToEntity() *entity.HealthcheckGRPC {
	if req == nil {
		return nil
	}
	return &entity.HealthcheckGRPC{
		Version:      req.Version,
		Addr:         req.Addr,
		Service:      req.Service,
		ReturnStatus: req.ReturnStatus,
	}
}

func (req *HealthcheckBaseReq) validate(_ string) []vld.Validator {
	// TODO: add validation
	return nil
}

func NewCreateHealthcheckReq() *CreateHealthcheckReq {
	return &CreateHealthcheckReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateHealthcheckReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateHealthcheckResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
