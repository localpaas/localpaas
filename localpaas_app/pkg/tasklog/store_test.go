package tasklog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_LocalAddAndGet(t *testing.T) {
	store := newStore("testkey", true, false, nil)

	f1 := NewInFrame("msg1", TsNow)
	f2 := NewOutFrame("msg2", TsNow)

	err := store.Add(context.Background(), f1, f2)
	assert.NoError(t, err)

	frames, err := store.GetData(context.Background(), 0)
	assert.NoError(t, err)
	assert.Len(t, frames, 2)
	assert.Equal(t, f1, frames[0])
	assert.Equal(t, f2, frames[1])
}
