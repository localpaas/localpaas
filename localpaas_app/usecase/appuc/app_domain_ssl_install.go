package appuc

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

const (
	certDirFileMode = 0755
)

func (uc *AppUC) InstallDomainSsl(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.InstallDomainSslReq,
) (*appdto.InstallDomainSslResp, error) {
	appData := &installSslData{}
	err := uc.loadAppDataForInstallDomainSsl(ctx, uc.db, req, appData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	saveDir := filepath.Join(config.Current.DataPathLetsEncryptEtc(), req.Domain)
	err = os.MkdirAll(saveDir, certDirFileMode)
	if err != nil {
		return nil, apperrors.New(err).WithMsgLog("failed to create directory to save certificates")
	}

	savePrivateKeyPath := filepath.Join(saveDir, "private.key")
	saveCertPath := filepath.Join(saveDir, "certificate.crt")
	_, err = uc.letsencryptClient.ObtainCertificate(ctx, []string{req.Domain}, savePrivateKeyPath, saveCertPath)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.InstallDomainSslResp{}, nil
}

type installSslData struct {
	App *entity.App
}

func (uc *AppUC) loadAppDataForInstallDomainSsl(
	ctx context.Context,
	db database.IDB,
	req *appdto.InstallDomainSslReq,
	data *installSslData,
) error {
	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if app.Status != base.AppStatusActive {
		return apperrors.Wrap(apperrors.ErrResourceInactive)
	}
	data.App = app

	return nil
}
