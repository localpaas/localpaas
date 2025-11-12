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
	SecretKey    string          `json:"secretKey"`
	AccessAction base.ActionType `json:"accessAction,omitempty"`
}

func (o *APIKey) IsHashed() bool {
	return strings.HasPrefix(o.SecretKey, base.SaltPrefix)
}

func (o *APIKey) Hash() error {
	if o.IsHashed() {
		return nil
	}
	secretHash, salt, err := randtoken.HashAsHex(o.SecretKey, base.DefaultSaltLen,
		apiKeyHashingKeyLen, apiKeyHashingIteration)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.SecretKey = cryptoutil.PackSecret(secretHash, salt)
	return nil
}

func (o *APIKey) MustHash() *APIKey {
	gofn.Must1(o.Hash())
	return o
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

func (s *Setting) ParseAPIKey() (*APIKey, error) {
	if s != nil && s.Data != "" && s.Type == base.SettingTypeAPIKey {
		res := &APIKey{}
		err := s.parseData(res)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		return res, nil
	}
	return nil, nil
}
