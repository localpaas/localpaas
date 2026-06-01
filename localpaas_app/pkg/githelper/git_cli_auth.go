package githelper

import (
	"os"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	sshKeyFileMode = 0600
)

func writeSshKeyFile(baseDir string, pemBytes []byte) (sshKeyFile string, err error) {
	fh, err := os.CreateTemp(baseDir, "git-ssh-*")
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	defer fh.Close()

	// NOTE: file will be removed along with the whole temp dir by the caller
	sshKeyFile = fh.Name()

	if err := os.Chmod(sshKeyFile, sshKeyFileMode); err != nil {
		return "", apperrors.Wrap(err)
	}

	if _, err := fh.Write(pemBytes); err != nil {
		return "", apperrors.Wrap(err)
	}

	if pemBytes[len(pemBytes)-1] != '\n' {
		if _, err := fh.Write([]byte("\n")); err != nil {
			return "", apperrors.Wrap(err)
		}
	}

	return sshKeyFile, nil
}
