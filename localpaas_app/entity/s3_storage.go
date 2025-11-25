package entity

import (
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type S3Storage struct {
	AccessKeyID string `json:"accessKeyId"`
	SecretKey   string `json:"secretKey"`
	Region      string `json:"region,omitempty"`
	Bucket      string `json:"bucket,omitempty"`
	Endpoint    string `json:"endpoint,omitempty"`
}

func (o *S3Storage) IsEncrypted() bool {
	return strings.HasPrefix(o.SecretKey, base.SaltPrefix)
}

func (o *S3Storage) Encrypt() (*S3Storage, error) {
	if o.IsEncrypted() {
		return o, nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.SecretKey, base.DefaultSaltLen)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.SecretKey = encrypted
	return o, nil
}

func (o *S3Storage) MustEncrypt() *S3Storage {
	return gofn.Must(o.Encrypt())
}

func (o *S3Storage) Decrypt() (*S3Storage, error) {
	if !o.IsEncrypted() {
		return o, nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.SecretKey)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.SecretKey = decrypted
	return o, nil
}

func (o *S3Storage) MustDecrypt() *S3Storage {
	return gofn.Must(o.Decrypt())
}

func (s *Setting) AsS3Storage() (*S3Storage, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*S3Storage)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &S3Storage{}
	if s.Data != "" && s.Type == base.SettingTypeS3Storage {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsS3Storage() *S3Storage {
	return gofn.Must(s.AsS3Storage())
}
