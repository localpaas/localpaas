// Source copied from: https://github.com/foomo/htpasswd

package htpasswd

import (
	"errors"
	"os"

	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

var (
	ErrPasswordRequired = errors.New("password is required")
)

// HashedPasswords name => hash
type HashedPasswords map[string]string

// HashAlgorithm enum for hashing algorithms
type HashAlgorithm string

const (
	// HashBCrypt bcrypt - recommended
	HashBCrypt = "bcrypt"
)

const (
	// PasswordSeparator separates passwords from hashes
	PasswordSeparator = ":"
	// LineSeparator separates password records
	LineSeparator = "\n"
)

const (
	fileMode = 0644
)

// Bytes bytes representation
func (hp HashedPasswords) Bytes() (passwordBytes []byte) {
	passwordBytes = []byte{}
	for name, hash := range hp {
		passwordBytes = append(passwordBytes, []byte(name+PasswordSeparator+hash+LineSeparator)...)
	}
	return passwordBytes
}

// WriteToFile put them to a file will be overwritten or created
func (hp HashedPasswords) WriteToFile(file string) error {
	err := os.WriteFile(file, hp.Bytes(), fileMode)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// SetPassword set a password for a user with a hashing algo
func (hp HashedPasswords) SetPassword(name, password string, hashAlgorithm HashAlgorithm) (err error) {
	if len(password) == 0 {
		return tracerr.Wrap(ErrPasswordRequired)
	}
	hash := ""
	switch hashAlgorithm { //nolint
	case HashBCrypt:
		hash, err = hashBcrypt(password)
	}
	if err != nil {
		return tracerr.Wrap(err)
	}
	hp[name] = hash
	return nil
}
