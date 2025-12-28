package docker

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/moby/moby/api/types/jsonstream"
)

type JSONMsg struct {
	*jsonstream.Message
}

func StartJSONMsgScanning(ctx context.Context, reader io.ReadCloser) (
	msgChan <-chan *JSONMsg, closeFunc func() error) {
	channel := make(chan *JSONMsg, 100) //nolint:mnd
	_, hasDeadline := ctx.Deadline()
	if hasDeadline {
		context.AfterFunc(ctx, func() {
			_ = reader.Close()
		})
	}

	go func() {
		defer close(channel)

		// Handle panic
		defer func() {
			_ = recover()
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
			select {
			case <-ctx.Done():
				return
			case channel <- &JSONMsg{Message: &jm}:
			}
		}
	}()

	return channel, func() error { return reader.Close() }
}

func StartJSONMsgBatchScanning(ctx context.Context, reader io.ReadCloser,
	maxSendingPeriod time.Duration, maxMsg int) (msgChan <-chan []*JSONMsg, closeFunc func() error) {
	msgBatchChan := make(chan []*JSONMsg, max(20, maxMsg)) //nolint:mnd

	go func() {
		_, hasDeadline := ctx.Deadline()
		if hasDeadline {
			context.AfterFunc(ctx, func() { _ = reader.Close() })
		}

		defer close(msgBatchChan)

		// Handle panic
		defer func() {
			_ = recover()
		}()

		// Close logs stream
		defer reader.Close()

		mu := sync.Mutex{}
		msgList := make([]*JSONMsg, 0, max(10, maxMsg)) //nolint:mnd
		var sendTimer *time.Timer

		sendFunc := func() {
			mu.Lock()
			if len(msgList) > 0 {
				msgBatchChan <- msgList
				msgList = msgList[:0]
			}
			sendTimer = nil
			mu.Unlock()
		}

		decoder := json.NewDecoder(reader)
		for {
			var jm jsonstream.Message
			err := decoder.Decode(&jm)
			if errors.Is(err, io.EOF) {
				break
			}
			msg := &JSONMsg{Message: &jm}

			sendImmediately := false
			mu.Lock()
			msgList = append(msgList, msg)
			if len(msgList) >= maxMsg { // Send data immediately
				sendImmediately = true
			} else if sendTimer == nil { // Send data in at most `period` duration
				sendTimer = time.AfterFunc(maxSendingPeriod, sendFunc)
			}
			mu.Unlock()

			if sendImmediately {
				sendFunc()
			}

			// Make sure to quit if the context is done
			select {
			case <-ctx.Done():
				return
			default:
				// Just continue
			}
		}
	}()

	return msgBatchChan, func() error { return reader.Close() }
}
