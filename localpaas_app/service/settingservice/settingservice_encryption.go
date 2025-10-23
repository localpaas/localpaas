package settingservice

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/pkg/reflectutil"
)

func (s *settingService) EncryptSetting(plaintext, salt []byte) ([]byte, error) {
	key := config.Current.App.Secret + string(salt)
	block, err := aes.NewCipher(reflectutil.UnsafeStrToBytes(key))
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, apperrors.Wrap(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (s *settingService) DecryptSetting(ciphertext, salt []byte) ([]byte, error) {
	key := config.Current.App.Secret + string(salt)
	block, err := aes.NewCipher(reflectutil.UnsafeStrToBytes(key))
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, apperrors.NewParamInvalid("ciphertext").
			WithMsgLog("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return plaintext, nil
}
