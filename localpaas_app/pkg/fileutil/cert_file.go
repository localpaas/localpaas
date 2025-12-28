package fileutil

import (
	"os"
	"path/filepath"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	certDirFileMode = 0o755
)

func WriteCerts(cert, privateKey []byte, saveDir, certFilePath, privateKeyFilePath string) (err error) {
	err = os.MkdirAll(saveDir, certDirFileMode)
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to create directory to save ssl certificate")
	}

	savePrivateKeyPath := filepath.Join(saveDir, privateKeyFilePath)
	saveCertPath := filepath.Join(saveDir, certFilePath)
	defer func() {
		if err == nil {
			return
		}
		// Try to remove all created files
		_ = os.Remove(savePrivateKeyPath)
		_ = os.Remove(saveCertPath)
	}()

	if savePrivateKeyPath != "" {
		privKeyFile, err := os.Create(savePrivateKeyPath)
		if err != nil {
			return apperrors.New(err).WithMsgLog("failed to create private key file")
		}
		_, err = privKeyFile.Write(privateKey)
		if err != nil {
			return apperrors.New(err).WithMsgLog("failed to write data to private key file")
		}
	}
	if saveCertPath != "" {
		certFile, err := os.Create(saveCertPath)
		if err != nil {
			return apperrors.New(err).WithMsgLog("failed to create certificate file")
		}
		_, err = certFile.Write(cert)
		if err != nil {
			return apperrors.New(err).WithMsgLog("failed to write data to certificate file")
		}
	}
	return nil
}
