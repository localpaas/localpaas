package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"

	"github.com/tiendc/gofn"
	"golang.org/x/crypto/argon2"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
)

const (
	keyLen        = 32
	hashIteration = 1
	hashMemory    = 64 * 1024
	hashThread    = 2
)

func makeKey(secret, salt []byte) []byte {
	return argon2.IDKey(secret, salt, hashIteration, hashMemory, hashThread, keyLen)
}

func EncryptEx(plaintext, salt, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(makeKey(secret, salt))
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

func Encrypt(plaintext, salt []byte) ([]byte, error) {
	return EncryptEx(plaintext, salt, []byte(config.Current.Secret))
}

// EncryptBase64 this encrypts the input and returns a string in form: `lpsalt:<salt> <secret>`
func EncryptBase64(plaintext string, saltLen int) (ciphertext string, err error) {
	if saltLen <= 0 {
		return plaintext, nil
	}

	saltBytes := gofn.RandToken(saltLen)
	ciphertextBytes, err := Encrypt([]byte(plaintext), saltBytes)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	ciphertext = base64.StdEncoding.EncodeToString(ciphertextBytes)
	salt := base64.StdEncoding.EncodeToString(saltBytes)
	return PackSecret(ciphertext, salt), nil
}

func DecryptEx(ciphertext, salt, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(makeKey(secret, salt))
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

func Decrypt(ciphertext, salt []byte) ([]byte, error) {
	return DecryptEx(ciphertext, salt, []byte(config.Current.Secret))
}

// DecryptBase64 this decrypts the input in form: `lpsalt:<salt> <secret>`
func DecryptBase64(secret string) (plaintext string, err error) {
	ciphertext, salt := UnpackSecret(secret)
	if salt == "" {
		return ciphertext, nil
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", apperrors.Wrap(err)
	}

	plaintextBytes, err := Decrypt(ciphertextBytes, saltBytes)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	plaintext = string(plaintextBytes)

	return plaintext, nil
}

func PackSecret(secret, salt string) string {
	return base.EncryptionSaltPrefix + salt + " " + secret
}

func UnpackSecret(secretText string) (secret string, salt string) {
	if !strings.HasPrefix(secretText, base.EncryptionSaltPrefix) {
		return secretText, ""
	}
	parts := strings.SplitN(secretText, " ", 2) //nolint:mnd
	if len(parts) != 2 {                        //nolint:mnd
		return secretText, ""
	}
	return parts[1], strings.TrimPrefix(parts[0], base.EncryptionSaltPrefix)
}
