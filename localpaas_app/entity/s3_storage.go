package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type S3Storage struct {
	AccessKeyID string `json:"accessKeyId"`
	SecretKey   string `json:"secretKey"`
	Salt        string `json:"salt,omitempty"`
	Region      string `json:"region,omitempty"`
	Bucket      string `json:"bucket,omitempty"`
}

func (o *S3Storage) IsEncrypted() bool {
	return o.Salt != ""
}

func (o *S3Storage) Encrypt() error {
	if o.Salt != "" {
		return nil
	}
	cipher, salt, err := cryptoutil.EncryptBase64(o.SecretKey, defaultSaltLen)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.SecretKey = cipher
	o.Salt = salt
	return nil
}

func (o *S3Storage) Decrypt() error {
	if o.Salt == "" {
		return nil
	}
	plain, err := cryptoutil.DecryptBase64(o.SecretKey, o.Salt)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.SecretKey = plain
	o.Salt = ""
	return nil
}

func (s *Setting) ParseS3Storage(decrypt bool) (*S3Storage, error) {
	if s != nil && s.Data != "" {
		res := &S3Storage{}
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
	return nil, nil
}
