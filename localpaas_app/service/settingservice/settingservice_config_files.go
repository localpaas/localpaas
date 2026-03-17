package settingservice

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/fileutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/htpasswd"
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
	certDir := config.Current.DataPathCerts()
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
	certDir := config.Current.DataPathCerts()
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

func (s *settingService) PersistBasicAuthConfigFiles(
	forceRecreate bool,
	settings ...*entity.Setting,
) error {
	if len(settings) == 0 {
		return nil
	}
	basicAuthDir := config.Current.DataPathNginxShareBasicAuth()
	err := os.MkdirAll(basicAuthDir, basicAuthDirFileMode)
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to create directory to save basic auth files")
	}

	for _, setting := range settings {
		authFile := filepath.Join(basicAuthDir, setting.ID)
		authFileExists, _ := fileutil.FileExists(authFile, true)

		if !forceRecreate && authFileExists {
			continue
		}
		basicAuth := setting.MustAsBasicAuth()

		hashedPasswd := htpasswd.HashedPasswords{}
		err = hashedPasswd.SetPassword(basicAuth.Username, basicAuth.Password.MustGetPlain(), htpasswd.HashBCrypt)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = hashedPasswd.WriteToFile(authFile)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

func (s *settingService) DeleteBasicAuthConfigFiles(
	settings ...*entity.Setting,
) error {
	if len(settings) == 0 {
		return nil
	}
	basicAuthDir := config.Current.DataPathNginxShareBasicAuth()
	for _, setting := range settings {
		err := os.Remove(filepath.Join(basicAuthDir, setting.ID))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return apperrors.Wrap(err)
		}
	}
	return nil
}
