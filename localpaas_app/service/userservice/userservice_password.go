package userservice

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"unicode"

	"github.com/tiendc/gofn"
	"golang.org/x/crypto/argon2"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

const (
	SkipCheckingCurrentPassword = ""
)

const (
	saltLength       = 16
	hashingIteration = 1
	hashingMemory    = 64 * 1024 // 64MB
	hashingThreads   = 4
	hashingKeyLength = 32
)

const (
	passwordMinLen             = 8
	passwordMaxLen             = 32
	passwordMustHaveLowercases = 1
	passwordMustHaveUppercases = 1
	passwordMustHaveDigits     = 1
	passwordMustHaveSymbols    = 1
)

var (
	errNotMeetRequirementParams = map[string]any{
		"MinLen":   passwordMinLen,
		"MaxLen":   passwordMaxLen,
		"Lowers":   passwordMustHaveLowercases,
		"Uppers":   passwordMustHaveUppercases,
		"Digits":   passwordMustHaveDigits,
		"Specials": passwordMustHaveSymbols,
	}

	specialCharset  = []rune("!@#$%^&*()_+-=[]{}|;':\",./<>?")
	mapSpecialChars = gofn.MapSliceToMap(specialCharset, func(k rune) (rune, struct{}) {
		return k, struct{}{}
	})
)

// ChangePassword updates user password with the new one.
// Action is rejected if `currPassword` does not match.
// Checking for current password is skipped if pass empty string.
func (s *userService) ChangePassword(user *entity.User, newPassword, currPassword string) error {
	if currPassword != "" {
		if err := s.VerifyPassword(user, currPassword); err != nil {
			return err
		}
	}

	// Verify password strength
	if err := s.CheckPasswordStrength(newPassword); err != nil {
		return err
	}

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash the password using Argon2 with recommended configuration
	hashedPass := argon2.IDKey([]byte(newPassword), salt, hashingIteration, hashingMemory,
		hashingThreads, hashingKeyLength)

	user.Password = hashedPass
	user.PasswordSalt = salt
	return nil
}

// VerifyPassword verifies the password matching the user data
func (s *userService) VerifyPassword(user *entity.User, password string) error {
	// We don't allow empty password
	if password == "" || len(user.Password) == 0 {
		return apperrors.New(apperrors.ErrPasswordMismatched)
	}
	passHash := argon2.IDKey([]byte(password), user.PasswordSalt, hashingIteration, hashingMemory,
		hashingThreads, hashingKeyLength)
	if !bytes.Equal(passHash, user.Password) {
		return apperrors.New(apperrors.ErrPasswordMismatched)
	}
	return nil
}

// CheckPasswordStrength checks the password if it meets the strength requirements
// TODO: consider checking password must not contain first/last name, email
// TODO: consider checking password must not be the same as last 3 history passwords
func (s *userService) CheckPasswordStrength(password string) error {
	chars := []rune(password)
	if len(chars) < passwordMinLen || len(chars) > passwordMaxLen {
		return apperrors.New(apperrors.ErrPasswordNotMeetRequirements).
			WithParams(errNotMeetRequirementParams).
			WithMsgLog("incorrect length: %d", len(chars))
	}
	lowers := 0
	uppers := 0
	digits := 0
	specials := 0
	for _, r := range chars {
		switch {
		case unicode.IsDigit(r):
			digits++
		case gofn.MapContainKeys(mapSpecialChars, r):
			specials++
		case unicode.IsLower(r):
			lowers++
		default:
			uppers++
		}
	}
	if lowers < passwordMustHaveLowercases || uppers < passwordMustHaveUppercases ||
		digits < passwordMustHaveDigits || specials < passwordMustHaveSymbols {
		return apperrors.New(apperrors.ErrPasswordNotMeetRequirements).
			WithParams(errNotMeetRequirementParams).
			WithMsgLog("lowers %d, uppers %d, digits %d, specials %d", lowers, uppers, digits, specials)
	}
	return nil
}
