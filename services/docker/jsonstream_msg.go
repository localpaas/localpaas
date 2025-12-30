package docker

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/moby/moby/api/types/jsonstream"

	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
)

type JSONMsg struct {
	*jsonstream.Message
}

func StartScanningJSONMsg(
	ctx context.Context,
	reader io.ReadCloser,
	options batchrecvchan.Options, // if zero, scan one by one
) (msgChan <-chan []*JSONMsg, closeFunc func() error) {
	batchChan := batchrecvchan.NewChan[*JSONMsg](options)

	_, hasDeadline := ctx.Deadline()
	if hasDeadline {
		context.AfterFunc(ctx, func() { _ = reader.Close() })
	}

	go func() {
		defer func() {
			_ = recover()
			batchChan.Close()
		}()

		// Close logs stream
		defer reader.Close()

		decoder := json.NewDecoder(reader)
		for {
			var jm jsonstream.Message
			err := decoder.Decode(&jm)
			if errors.Is(err, io.EOF) {
				break
			}
			msg := &JSONMsg{Message: &jm}

			select {
			case <-ctx.Done(): // Make sure to quit if the context is done
				return
			default:
				batchChan.Send(msg)
			}
		}
	}()

	return batchChan.Receiver(), func() error { return reader.Close() }
}
