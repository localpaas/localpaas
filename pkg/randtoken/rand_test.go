package randtoken

import (
	"testing"

	"github.com/stretchr/testify/assert"
	strvld "github.com/tiendc/go-validator/base/string"
)

func Test_New(t *testing.T) {
	t.Run("token length check", func(t *testing.T) {
		token, err := New(32)
		assert.Nil(t, err)
		assert.Equal(t, 32*2, len(token)) // Token is in hex form
		isHex, _ := strvld.IsHexadecimal(token)
		assert.True(t, isHex)
	})

	t.Run("tokens must differ from each other", func(t *testing.T) {
		token1, _ := New(16)
		token2, _ := New(16)
		assert.NotEqual(t, token1, token2)
	})
}
