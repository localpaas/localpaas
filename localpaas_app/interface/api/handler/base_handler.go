package handler

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/pkg/httputil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/strutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/translation"
)

var (
	defaultParseFuncs = []mapstructure.DecodeHookFunc{
		// TODO: add parse hook functions for specific types

		// This helps parse timeutil.Date from string value
		timeutil.MapstructureParseDateFunc(),

		// NOTE: Parse slice should be the last hook function as some
		// types may be treated as slice, such as UUID can be []byte.
		mapstructure.StringToSliceHookFunc(","),
	}
)

// BaseHandler base handler for every other handler to inherit from
type BaseHandler struct{}

func (h *BaseHandler) RequestCtx(ginCtx *gin.Context) context.Context {
	// TODO: for now just use the input context without any conversion
	return ginCtx
}

// RenderResponse renders response to client as JSON
func (h *BaseHandler) RenderResponse(ctx *gin.Context, status int, body any) {
	ctx.JSON(status, body)
}

// RenderError renders errors to client as JSON
func (h *BaseHandler) RenderError(ctx *gin.Context, err error) {
	// Parse the error
	errInfo, errType := apperrors.ParseError(err, h.ParseRequestLang(ctx))
	h.NotifyError(ctx, err, errInfo, errType)

	// Remove the error details from the response if we are in production env
	if config.Current.IsProdEnv() {
		errInfo.Cause = ""
		errInfo.DebugLog = ""

		// Also remove the details from InnerErrors
		for i := range errInfo.InnerErrors {
			errInfo.InnerErrors[i].Cause = ""
		}
	}

	// Always remove the stack trace regardless of env due to concern of response body length
	errInfo.StackTrace = ""

	ctx.JSON(errInfo.Status, errInfo)
}

// NotifyError notifies error
func (h *BaseHandler) NotifyError(_ *gin.Context, _ error, _ *apperrors.ErrorInfo,
	_ apperrors.ErrLevel) {
	// TODO: do nothing now
}

// parsePagination parses paging and sorting params
//
//nolint:gocognit
func (h *BaseHandler) parsePagination(ctx *gin.Context, paging *basedto.Paging) error {
	if paging.Limit == 0 {
		paging.Limit = basedto.PageLimitDefault
	}

	if limitStr := ctx.Query("pageLimit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > basedto.PageLimitMax {
			return apperrors.NewParamInvalid("pageLimit")
		}
		paging.Limit = limit
	}
	if offsetStr := ctx.Query("pageOffset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return apperrors.NewParamInvalid("pageOffset")
		}
		if offset > 0 {
			paging.Offset = offset
		}
	}

	//nolint:nestif
	if sortQuery := ctx.Query("sort"); sortQuery != "" {
		orders := basedto.Orders{}
		sort := strings.Split(sortQuery, ",")
		for _, str := range sort {
			if strings.HasPrefix(str, "-") {
				if len(str) == 1 {
					return apperrors.NewParamInvalid("sort")
				}
				orders.Add(&basedto.Order{Direction: basedto.DirectionDesc, ColumnName: str[1:]})
			} else {
				orders.Add(&basedto.Order{Direction: basedto.DirectionAsc, ColumnName: str})
			}
		}
		if len(paging.Sort) > 0 {
			orders.Add(paging.Sort...)
		}
		paging.Sort = orders
	}
	// Converts camelCase to snake_case
	for _, order := range paging.Sort {
		// Column name may include table such as `task.created_at`
		parts := strings.Split(order.ColumnName, ".")
		for i := range parts {
			parts[i] = strutil.ToSnakeCase(parts[i])
		}
		order.ColumnName = strings.Join(parts, ".")
	}

	return nil
}

// parseQuery parse request query from URL of style /api?key1=val1&key2=val2
func (h *BaseHandler) parseQuery(ctx *gin.Context, query any) error {
	mapQuery := map[string]string{}
	urlQueries := ctx.Request.URL.Query()
	for k, v := range urlQueries {
		mapQuery[k] = strings.Join(v, ",")
	}
	// NOTE: caller should call ParseMultipartForm()/ParseForm() if the request contains form data.
	for formParam, value := range ctx.Request.Form {
		mapQuery[formParam] = strings.Join(value, ",")
	}
	if len(mapQuery) == 0 {
		return nil
	}

	config := &mapstructure.DecoderConfig{
		Result:           query,
		WeaklyTypedInput: true,
	}

	if len(defaultParseFuncs) == 1 {
		config.DecodeHook = defaultParseFuncs[0]
	} else if len(defaultParseFuncs) > 1 {
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(defaultParseFuncs...)
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return apperrors.New(err)
	}
	if err = decoder.Decode(mapQuery); err != nil {
		return apperrors.NewParamInvalid("query").WithCause(err)
	}

	return nil
}

// ParseRequest parses request query, then validate the request.
//
// The request struct should be defined similarly as:
//
//	type ListProjectReq struct {
//	   Key1 int    `mapstructure:"key1"`
//	   Key2 string `mapstructure:"key2"`
//	}
func (h *BaseHandler) ParseRequest(ctx *gin.Context, reqStruct any, paging *basedto.Paging) error {
	if err := h.parseQuery(ctx, reqStruct); err != nil {
		return err
	}
	if paging != nil {
		if err := h.parsePagination(ctx, paging); err != nil {
			return err
		}
	}

	// Execute custom modifier for the request input
	if modifier, ok := reqStruct.(basedto.ReqModifier); ok {
		if err := modifier.ModifyRequest(); err != nil {
			return apperrors.New(err)
		}
	}

	// Execute custom validator for the request input
	if validator, ok := reqStruct.(basedto.ReqValidator); ok {
		if vldErrs := validator.Validate(); len(vldErrs) > 0 {
			return vldErrs
		}
	}

	return nil
}

// ParseStringParam parse string from context param
func (h *BaseHandler) ParseStringParam(ctx *gin.Context, paramName string) (string, error) {
	idStr := ctx.Params.ByName(paramName)
	if idStr == "" {
		return "", apperrors.New(apperrors.ErrBadRequest).
			WithMsgLog("require param `%s` of type string in URL", paramName)
	}
	return idStr, nil
}

// ParseIntParam parse uint from context param
func (h *BaseHandler) ParseIntParam(ctx *gin.Context, paramName string) (int, error) {
	value, err := strconv.ParseInt(ctx.Param(paramName), 10, 64) //nolint:mnd
	if err != nil {
		return 0, apperrors.New(apperrors.ErrBadRequest).WithCause(err).
			WithMsgLog("require param `%s` of type int in URL", paramName)
	}
	return int(value), nil
}

// ParseUintParam parse uint from context param
func (h *BaseHandler) ParseUintParam(ctx *gin.Context, paramName string) (uint, error) {
	value, err := strconv.ParseUint(ctx.Param(paramName), 10, 64) //nolint:mnd
	if err != nil {
		return 0, apperrors.New(apperrors.ErrBadRequest).WithCause(err).
			WithMsgLog("require param `%s` of type uint in URL", paramName)
	}
	return uint(value), nil
}

// ParseJSONBody parse request body as JSON, then validate the result
func (h *BaseHandler) ParseJSONBody(ctx *gin.Context, reqStruct any) error {
	var buf bytes.Buffer
	ctx.Request.Body = io.NopCloser(io.TeeReader(ctx.Request.Body, &buf))

	if err := ctx.ShouldBindJSON(reqStruct); err != nil && buf.Len() > 0 {
		if handler, ok := reqStruct.(basedto.ReqParsingErrorHandler); ok {
			if newErr := handler.HandleParsingError(err); newErr != nil {
				return apperrors.New(newErr)
			}
		}
		return apperrors.New(apperrors.ErrBadRequest).WithCause(err)
	}

	ctx.Request.Body = io.NopCloser(&buf)

	// Execute custom modifier for the request input
	if modifier, ok := reqStruct.(basedto.ReqModifier); ok {
		if err := modifier.ModifyRequest(); err != nil {
			return apperrors.New(err)
		}
	}

	// Execute custom validator for the request input
	if validator, ok := reqStruct.(basedto.ReqValidator); ok {
		if vldErrs := validator.Validate(); len(vldErrs) > 0 {
			return vldErrs
		}
	}

	return nil
}

func (h *BaseHandler) ParseRequestLang(ctx *gin.Context) translation.Lang {
	return httputil.ParseRequestLang(ctx.GetHeader("Accept-Language"))
}
