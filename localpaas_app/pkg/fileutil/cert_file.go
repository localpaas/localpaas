package fileutil

import (
	"os"
	"path/filepath"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	certDirFileMode = 0o755
	certFileMode    = 0o644
)

func WriteCerts(cert, privateKey []byte, saveDir, certFile, keyFile string, overwrite bool) (err error) {
	err = os.MkdirAll(saveDir, certDirFileMode)
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to create directory to save ssl certificate")
	}

	certPath := filepath.Join(saveDir, certFile)
	keyPath := filepath.Join(saveDir, keyFile)

	defer func() {
		if err == nil {
			return
		}
		// Try to remove all created files
		_ = os.Remove(keyPath)
		_ = os.Remove(certPath)
	}()

	if certPath != "" {
		if overwrite || !gofn.Head(FileExists(certPath, true)) {
			err = os.WriteFile(certPath, cert, certFileMode)
			if err != nil {
				return apperrors.New(err).WithMsgLog("failed to write cert file to %s", certPath)
			}
		}
	}
	if keyPath != "" {
		if overwrite || !gofn.Head(FileExists(keyPath, true)) {
			err = os.WriteFile(keyPath, privateKey, certFileMode)
			if err != nil {
				return apperrors.New(err).WithMsgLog("failed to write key file to %s", keyPath)
			}
		}
	}

	return nil
}
