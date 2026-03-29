package settingservice

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/fileutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

const (
	certDirFileMode      = 0o755
	basicAuthDirFileMode = 0o755
)

func (s *settingService) PersistSSLConfigFiles(
	forceRecreate bool,
	settings ...*entity.Setting,
) error {
	if len(settings) == 0 {
		return nil
	}
	certDir := config.Current.DataPathSslCerts()
	err := os.MkdirAll(certDir, certDirFileMode)
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to create directory to save cert files")
	}

	for _, setting := range settings {
		certFile := setting.ID + ".crt"
		keyFile := setting.ID + ".key"
		certFileExists, _ := fileutil.FileExists(filepath.Join(certDir, certFile), true)
		keyFileExists, _ := fileutil.FileExists(filepath.Join(certDir, keyFile), true)

		if !forceRecreate && certFileExists && keyFileExists {
			continue
		}

		ssl := setting.MustAsSSL()
		certBytes := reflectutil.UnsafeStrToBytes(ssl.Certificate)
		keyBytes := reflectutil.UnsafeStrToBytes(ssl.PrivateKey.MustGetPlain())

		err := fileutil.WriteCerts(certBytes, keyBytes, certDir, certFile, keyFile, true)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

func (s *settingService) DeleteSSLConfigFiles(
	settings ...*entity.Setting,
) error {
	if len(settings) == 0 {
		return nil
	}
	certDir := config.Current.DataPathSslCerts()
	for _, setting := range settings {
		certFile := setting.ID + ".crt"
		keyFile := setting.ID + ".key"

		err := os.Remove(filepath.Join(certDir, certFile))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return apperrors.Wrap(err)
		}

		err = os.Remove(filepath.Join(certDir, keyFile))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return apperrors.Wrap(err)
		}
	}
	return nil
}
