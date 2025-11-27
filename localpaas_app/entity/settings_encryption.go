package entity

import (
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

type EncryptedField struct {
	encrypted string
	decrypted string
}

func (s *EncryptedField) MarshalJSON() ([]byte, error) {
	encrypted, err := s.encrypt()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return reflectutil.UnsafeStrToBytes(gofn.StringWrap(encrypted, "\"")), nil
}

func (s *EncryptedField) UnmarshalJSON(data []byte) error {
	s.encrypted = gofn.StringUnwrap(reflectutil.UnsafeBytesToStr(data), "\"")
	return nil
}

func (s *EncryptedField) String() string {
	if s.decrypted != "" {
		return s.decrypted
	}
	return s.encrypted
}

func (s *EncryptedField) IsEncrypted() bool {
	return s.encrypted != "" && s.decrypted == ""
}

func (s *EncryptedField) GetPlain() (string, error) {
	decrypted, err := s.decrypt()
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return decrypted, nil
}

func (s *EncryptedField) MustGetPlain() string {
	return gofn.Must(s.GetPlain())
}

func (s *EncryptedField) GetEncrypted() (string, error) {
	encrypted, err := s.encrypt()
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return encrypted, nil
}

func (s *EncryptedField) MustGetEncrypted() string {
	return gofn.Must(s.GetEncrypted())
}

func (s *EncryptedField) Set(value string) {
	if strings.HasPrefix(value, base.SaltPrefix) {
		s.encrypted = value
		s.decrypted = ""
	} else {
		s.decrypted = value
		s.encrypted = ""
	}
}

func (s *EncryptedField) encrypt() (string, error) {
	if s.encrypted != "" {
		return s.encrypted, nil
	}
	encrypted, err := cryptoutil.EncryptBase64(s.decrypted, base.DefaultSaltLen)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	s.encrypted = encrypted
	return encrypted, nil
}

func (s *EncryptedField) decrypt() (string, error) {
	if s.decrypted != "" {
		return s.decrypted, nil
	}
	decrypted, err := cryptoutil.DecryptBase64(s.encrypted)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	s.decrypted = decrypted
	return decrypted, nil
}

func NewEncryptedField(value string) EncryptedField {
	resp := EncryptedField{}
	resp.Set(value)
	return resp
}
