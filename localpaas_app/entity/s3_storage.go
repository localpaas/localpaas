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

func (o *S3Storage) Encrypt() error {
	if o.IsEncrypted() {
		return nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.SecretKey, base.DefaultSaltLen)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.SecretKey = encrypted
	return nil
}

func (o *S3Storage) MustEncrypt() *S3Storage {
	gofn.Must1(o.Encrypt())
	return o
}

func (o *S3Storage) Decrypt() error {
	if !o.IsEncrypted() {
		return nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.SecretKey)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.SecretKey = decrypted
	return nil
}

func (s *Setting) ParseS3Storage(decrypt bool) (*S3Storage, error) {
	res := &S3Storage{}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeS3Storage {
		err := s.parseData(res)
		if err != nil {
			return nil, err
		}
		if decrypt {
			if err = res.Decrypt(); err != nil {
				return nil, apperrors.Wrap(err)
			}
		}
		return res, nil
	}
	return res, nil
}
