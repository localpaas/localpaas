package apikeydto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

const (
	expirationYearMax = 1
)

type CreateAPIKeyReq struct {
	AccessAction base.ActionType `json:"accessAction"`
	Expiration   time.Time       `json:"expiration"`
}

func NewCreateAPIKeyReq() *CreateAPIKeyReq {
	return &CreateAPIKeyReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateAPIKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	timeNow := timeutil.NowUTC()
	validators = append(validators, basedto.ValidateStrIn(&req.AccessAction, false, base.AllActionTypes,
		"accessAction")...)
	validators = append(validators, basedto.ValidateTime(&req.Expiration, false, timeNow,
		timeNow.AddDate(expirationYearMax, 0, 0), "expiration")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateAPIKeyResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *APIKeyDataResp   `json:"data"`
}

type APIKeyDataResp struct {
	ID        string `json:"id"`
	KeyID     string `json:"keyId"`
	SecretKey string `json:"secretKey"`
}
