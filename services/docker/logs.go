package docker

import (
	"bufio"
	"context"
	"io"
	"sync"
	"time"
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

func StartLogScanning(ctx context.Context, logsReader io.ReadCloser) chan *LogFrame {
	scanner := bufio.NewScanner(logsReader)
	channel := make(chan *LogFrame, 100) //nolint:mnd
	_, hasDeadline := ctx.Deadline()
	if hasDeadline {
		context.AfterFunc(ctx, func() {
			_ = logsReader.Close()
		})
	}

	go func() {
		defer close(channel)

		// Handle panic
		defer func() {
			_ = recover()
		}()

		// Close logs stream
		defer logsReader.Close()

		for scanner.Scan() {
			// Format structure of the logs data, see:
			// https://docs.docker.com/reference/api/engine/version/v1.51/#tag/Container/operation/ContainerAttach
			logBytes := scanner.Bytes()
			var logType LogType

			//nolint:mnd
			if len(logBytes) > 8 {
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

			frame := &LogFrame{
				Data: string(logBytes),
				Type: logType,
			}
			select {
			case <-ctx.Done():
				return
			case channel <- frame:
			}
		}
	}()

	return channel
}

func StartLogBatchScanning(ctx context.Context, logsReader io.ReadCloser, period time.Duration,
	maxFrame int) chan []*LogFrame {
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
			// Format structure of the logs data, see:
			// https://docs.docker.com/reference/api/engine/version/v1.51/#tag/Container/operation/ContainerAttach
			logBytes := scanner.Bytes()
			var logType LogType

			//nolint:mnd
			if len(logBytes) > 8 {
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

			sendImmediately := false
			mu.Lock()
			logFrames = append(logFrames, &LogFrame{
				Data: string(logBytes),
				Type: logType,
			})
			if len(logFrames) == maxFrame { // Send data immediately
				sendImmediately = true
			} else if sendTimer == nil { // Send data in at most `period` duration
				sendTimer = time.AfterFunc(period, sendFunc)
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

	return logBatchChan
}
