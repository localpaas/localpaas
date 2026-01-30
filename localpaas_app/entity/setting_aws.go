package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentAWSVersion = 1
)

type AWS struct {
	AccessKeyID string         `json:"accessKeyId"`
	SecretKey   EncryptedField `json:"secretKey"`
	Region      string         `json:"region,omitempty"`
}

func (s *AWS) GetType() base.SettingType {
	return base.SettingTypeAWS
}

func (s *AWS) GetRefSettingIDs() []string {
	return nil
}

func (s *AWS) MustDecrypt() *AWS {
	s.SecretKey.MustGetPlain()
	return s
}

func (s *Setting) AsAWS() (*AWS, error) {
	return parseSettingAs(s, func() *AWS { return &AWS{} })
}

func (s *Setting) MustAsAWS() *AWS {
	return gofn.Must(s.AsAWS())
}
