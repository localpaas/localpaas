package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type S3Storage struct {
	AccessKeyID string         `json:"accessKeyId"`
	SecretKey   EncryptedField `json:"secretKey"`
	Region      string         `json:"region,omitempty"`
	Bucket      string         `json:"bucket,omitempty"`
	Endpoint    string         `json:"endpoint,omitempty"`
}

func (o *S3Storage) MustDecrypt() *S3Storage {
	o.SecretKey.MustGetPlain()
	return o
}

func (s *Setting) AsS3Storage() (*S3Storage, error) {
	return parseSettingAs(s, base.SettingTypeS3Storage, func() *S3Storage { return &S3Storage{} })
}

func (s *Setting) MustAsS3Storage() *S3Storage {
	return gofn.Must(s.AsS3Storage())
}
