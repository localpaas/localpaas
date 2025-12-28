package docker

import (
	"bufio"
	"context"
	"io"
	"sync"
	"time"
)

type LogKind string

const (
	LogKindContainer LogKind = "container"
)

type LogType string

const (
	LogTypeStdout LogType = "out"
	LogTypeStdin  LogType = "in"
	LogTypeStderr LogType = "err"
)

type LogFrame struct {
	Data string  `json:"data"`
	Type LogType `json:"type"`
}

func StartLogScanning(ctx context.Context, reader io.ReadCloser, kind LogKind) (
	logChan <-chan *LogFrame, closeFunc func() error) {
	scanner := bufio.NewScanner(reader)
	channel := make(chan *LogFrame, 100) //nolint:mnd
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

		for scanner.Scan() {
			logFrame := parseLogFrame(scanner.Bytes(), kind)
			select {
			case <-ctx.Done():
				return
			case channel <- logFrame:
			}
		}
	}()

	return channel, func() error { return reader.Close() }
}

func StartLogBatchScanning(ctx context.Context, logsReader io.ReadCloser, kind LogKind,
	maxSendingPeriod time.Duration, maxFrame int) (logChan <-chan []*LogFrame, closeFunc func() error) {
	logBatchChan := make(chan []*LogFrame, max(20, maxFrame)) //nolint:mnd

	go func() {
		_, hasDeadline := ctx.Deadline()
		if hasDeadline {
			context.AfterFunc(ctx, func() { _ = logsReader.Close() })
		}

		defer close(logBatchChan)

		// Handle panic
		defer func() {
			_ = recover()
		}()

		// Close logs stream
		defer logsReader.Close()

		mu := sync.Mutex{}
		logFrames := make([]*LogFrame, 0, max(20, maxFrame)) //nolint:mnd
		var sendTimer *time.Timer

		sendFunc := func() {
			mu.Lock()
			if len(logFrames) > 0 {
				logBatchChan <- logFrames
				logFrames = logFrames[:0]
			}
			sendTimer = nil
			mu.Unlock()
		}

		scanner := bufio.NewScanner(logsReader)
		for scanner.Scan() {
			logFrame := parseLogFrame(scanner.Bytes(), kind)

			sendImmediately := false
			mu.Lock()
			logFrames = append(logFrames, logFrame)
			if len(logFrames) >= maxFrame { // Send data immediately
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

	return logBatchChan, func() error { return logsReader.Close() }
}

func parseLogFrame(logBytes []byte, kind LogKind) *LogFrame {
	var logType LogType
	// Format structure of the logs data, see:
	// https://docs.docker.com/reference/api/engine/version/v1.51/#tag/Container/operation/ContainerAttach
	//nolint:mnd
	if kind == LogKindContainer && len(logBytes) > 8 {
		switch logBytes[0] {
		case 0:
			logType = LogTypeStdin
		case 1:
			logType = LogTypeStdout
		case 2:
			logType = LogTypeStderr
		}
		logBytes = logBytes[8:]
	}
	return &LogFrame{
		Data: string(logBytes),
		Type: logType,
	}
}
