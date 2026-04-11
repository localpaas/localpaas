package appsettingsuc

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appsettingsuc/appsettingsdto"
)

const (
	defaultCheckPortTimeout = time.Second * 5
)

func (uc *UC) CheckAppContainerPort(
	ctx context.Context,
	auth *basedto.Auth,
	req *appsettingsdto.CheckAppContainerPortReq,
) (*appsettingsdto.CheckAppContainerPortResp, error) {
	app, err := uc.appRepo.GetByID(ctx, uc.db, req.ProjectID, req.AppID,
		bunex.SelectExcludeColumns(entity.AppDefaultExcludeColumns...),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	address := net.JoinHostPort(app.Key, strconv.Itoa(int(req.Port))) //nolint
	timeout := gofn.Coalesce(req.Timeout.ToDuration(), defaultCheckPortTimeout)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err == nil && conn != nil {
		defer conn.Close()
		return &appsettingsdto.CheckAppContainerPortResp{
			Data: &appsettingsdto.CheckAppContainerPortDataResp{Open: true},
		}, nil
	}

	return &appsettingsdto.CheckAppContainerPortResp{
		Data: &appsettingsdto.CheckAppContainerPortDataResp{Open: false},
	}, nil
}
