package entity

import (
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/randtoken"
)

const (
	apiKeyHashingKeyLen    = 32
	apiKeyHashingIteration = 1
)

type APIKey struct {
	KeyID        string              `json:"keyId"`
	SecretKey    string              `json:"secretKey"`
	AccessAction *base.AccessActions `json:"accessAction,omitempty"`
}

func (o *APIKey) IsHashed() bool {
	return strings.HasPrefix(o.SecretKey, base.SaltPrefix)
}

func (o *APIKey) Hash() (*APIKey, error) {
	if o.IsHashed() {
		return o, nil
	}
	secretHash, salt, err := randtoken.HashAsHex(o.SecretKey, base.DefaultSaltLen,
		apiKeyHashingKeyLen, apiKeyHashingIteration)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.SecretKey = cryptoutil.PackSecret(secretHash, salt)
	return o, nil
}

func (o *APIKey) MustHash() *APIKey {
	return gofn.Must(o.Hash())
}

func (o *APIKey) VerifyHash(secretKey string) error {
	hash, salt := cryptoutil.UnpackSecret(o.SecretKey)
	var matched bool
	if salt == "" {
		matched = hash == secretKey
	} else {
		matched = randtoken.VerifyHashHex(secretKey, hash, salt, apiKeyHashingKeyLen, apiKeyHashingIteration)
	}
	if !matched {
		return apperrors.Wrap(apperrors.ErrAPIKeyMismatched)
	}
	return nil
}

func (s *Setting) AsAPIKey() (*APIKey, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*APIKey)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &APIKey{}
	if s.Data != "" && s.Type == base.SettingTypeAPIKey {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsAPIKey() *APIKey {
	return gofn.Must(s.AsAPIKey())
}
