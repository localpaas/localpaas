package taskhealthcheck

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/httpclient"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

func (e *Executor) doHealthcheckREST(
	ctx context.Context,
	data *taskData,
) (err error) {
	healthchk := data.Healthcheck.REST
	if data.Output.REST == nil {
		data.Output.REST = &entity.TaskHealthcheckOutputREST{}
	}

	reqCtx := ctx
	if data.Healthcheck.Timeout > 0 {
		ctx, cancel := context.WithTimeout(ctx, data.Healthcheck.Timeout.ToDuration())
		defer cancel()
		reqCtx = ctx
	}

	method := gofn.Coalesce(healthchk.Method, "GET")
	req, err := http.NewRequestWithContext(reqCtx, method, healthchk.URL, nil)
	if err != nil {
		return apperrors.Wrap(err)
	}
	req.Header.Set("Content-Type", gofn.Coalesce(healthchk.ContentType, "application/json"))

	resp, err := httpclient.DefaultClient.Do(req)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	data.Output.REST.ReturnCode = resp.StatusCode
	if healthchk.ReturnCode > 0 && healthchk.ReturnCode != resp.StatusCode {
		return apperrors.Wrap(apperrors.ErrActionFailed)
	}

	//nolint:nestif
	if healthchk.ReturnText != "" || healthchk.ReturnJSON != "" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return apperrors.Wrap(err)
		}
		bodyStr := reflectutil.UnsafeBytesToStr(body)
		if healthchk.ReturnText != "" {
			data.Output.REST.ReturnText = bodyStr
		}
		if healthchk.ReturnJSON != "" {
			data.Output.REST.ReturnJSON = bodyStr
		}

		if healthchk.ReturnText != "" && healthchk.ReturnText != bodyStr {
			return apperrors.Wrap(apperrors.ErrActionFailed)
		}
		if healthchk.ReturnJSON != "" && !compareRespJSON(body, healthchk.ReturnJSON) {
			return apperrors.Wrap(apperrors.ErrActionFailed)
		}
	}

	return nil
}

func compareRespJSON(resp []byte, expected string) bool {
	var actualObj, expectedObj any
	err := json.Unmarshal(resp, &actualObj)
	if err != nil {
		return false
	}
	err = json.Unmarshal(reflectutil.UnsafeStrToBytes(expected), &expectedObj)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(actualObj, expectedObj)
}
