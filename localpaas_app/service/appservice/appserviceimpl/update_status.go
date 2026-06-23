package appserviceimpl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/moby/moby/api/types/swarm"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

const (
	labelLocalPaaSAppPrevServiceMode = "localpaas.app.prevServiceMode"
)

func (s *service) OnAppStatusChanged(ctx context.Context, app *entity.App, oldStatus base.AppStatus) error {
	if app.Status == oldStatus {
		return nil
	}
	if app.Status == base.AppStatusDisabled {
		return s.onAppDisabled(ctx, app)
	}
	if app.Status == base.AppStatusActive {
		return s.onAppEnabled(ctx, app)
	}
	return nil
}

func (s *service) onAppDisabled(ctx context.Context, app *entity.App) error {
	if app.ServiceID == "" {
		return nil
	}

	inspect, err := s.dockerManager.ServiceInspect(ctx, app.ServiceID)
	if err != nil {
		return apperrors.New(err)
	}
	service := &inspect.Service

	prevSvcMode, err := json.Marshal(service.Spec.Mode)
	if err != nil {
		return apperrors.New(err)
	}
	service.Spec.Labels[labelLocalPaaSAppPrevServiceMode] = string(prevSvcMode)

	// Scale down to 0
	service.Spec.Mode = swarm.ServiceMode{
		Replicated: &swarm.ReplicatedService{
			Replicas: new(uint64(0)),
		},
	}

	err = gofn.ExecRetry(func() error {
		_, err := s.dockerManager.ServiceUpdate(ctx, app.ServiceID, &service.Version, &service.Spec)
		return apperrors.New(err)
	}, 2, 5*time.Second) //nolint:mnd
	if err != nil {
		return apperrors.New(err)
	}

	return nil
}

func (s *service) onAppEnabled(ctx context.Context, app *entity.App) error {
	if app.ServiceID == "" {
		return nil
	}

	inspect, err := s.dockerManager.ServiceInspect(ctx, app.ServiceID)
	if err != nil {
		return apperrors.New(err)
	}
	service := &inspect.Service

	prevSvcModeStr := service.Spec.Labels[labelLocalPaaSAppPrevServiceMode]
	if prevSvcModeStr != "" {
		mode := swarm.ServiceMode{}
		err = json.Unmarshal(reflectutil.UnsafeStrToBytes(prevSvcModeStr), &mode)
		if err != nil {
			return apperrors.New(err)
		}
		service.Spec.Mode = mode
		delete(service.Spec.Labels, labelLocalPaaSAppPrevServiceMode)
	} else {
		service.Spec.Mode = swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: new(uint64(1)),
			},
		}
	}

	err = gofn.ExecRetry(func() error {
		_, err := s.dockerManager.ServiceUpdate(ctx, app.ServiceID, &service.Version, &service.Spec)
		return apperrors.New(err)
	}, 2, 5*time.Second) //nolint:mnd
	if err != nil {
		return apperrors.New(err)
	}

	return nil
}
