package tasklog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/localpaas/localpaas/localpaas_app/pkg/redact"
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

func TestStore_AddRedacted(t *testing.T) {
	r := redact.New([]string{"secret123"})

	store := newStore("testkey", true, false, nil)
	store.SetRedactor(r)

	f1 := NewInFrame("normal message", TsNow)
	f2 := NewOutFrame("leak secret123 in log", TsNow)

	err := store.AddRedacted(context.Background(), f1, f2)
	assert.NoError(t, err)

	frames, err := store.GetData(context.Background(), 0)
	assert.NoError(t, err)
	assert.Len(t, frames, 2)

	assert.Equal(t, "normal message", frames[0].Data)
	assert.Equal(t, "leak ******** in log", frames[1].Data)
}
