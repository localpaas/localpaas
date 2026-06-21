package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentAccessTokenVersion = 1
)

var _ = registerSettingParser(base.SettingTypeAccessToken, &accessTokenParser{})

type accessTokenParser struct {
}

func (s *accessTokenParser) New() SettingData {
	return &AccessToken{}
}

type AccessToken struct {
	User    string         `json:"user"`
	Token   EncryptedField `json:"token"`
	BaseURL string         `json:"baseURL"`
}

func (s *AccessToken) GetType() base.SettingType {
	return base.SettingTypeAccessToken
}

func (s *AccessToken) GetRefObjectIDs() *RefObjectIDs {
	return &RefObjectIDs{}
}

func (s *AccessToken) CalcResLinks(setting *Setting) []*ResLink {
	return s.GetRefObjectIDs().CalcResLinks(base.ResourceTypeSetting, setting.ID)
}

func (s *AccessToken) Migrate(setting *Setting) (hasChange bool, err error) {
	if setting.Version == CurrentAccessTokenVersion {
		return false, nil
	}
	if setting.Version > CurrentAccessTokenVersion {
		return false, apperrors.New(apperrors.ErrDataVerNewerThanSystemVer)
	}

	// TODO: add migration if we make any change

	setting.Version = CurrentAccessTokenVersion
	setting.UpdateVer++
	setting.MustSetData(s)
	return true, nil
}

func (s *AccessToken) MustDecrypt() *AccessToken {
	s.Token.MustGetPlain()
	return s
}

func (s *Setting) AsAccessToken() (*AccessToken, error) {
	return parseSettingAs[*AccessToken](s)
}

func (s *Setting) MustAsAccessToken() *AccessToken {
	return gofn.Must(s.AsAccessToken())
}
