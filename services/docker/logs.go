package docker

import (
	"bufio"
	"context"
	"io"

	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
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

func StartScanningLog(
	ctx context.Context,
	reader io.ReadCloser,
	options batchrecvchan.Options, // if zero, scan one by one
) (logChan <-chan []*LogFrame, closeFunc func() error) {
	batchChan := batchrecvchan.NewChan[*LogFrame](options)

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

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			logFrame := parseLogFrame(scanner.Bytes())
			select {
			case <-ctx.Done(): // Make sure to quit if the context is done
				return
			default:
				batchChan.Send(logFrame)
			}
		}
	}()

	return batchChan.Receiver(), func() error { return reader.Close() }
}

func parseLogFrame(logBytes []byte) *LogFrame {
	var logType LogType
	// Format structure of the logs data, see:
	// https://docs.docker.com/reference/api/engine/version/v1.51/#tag/Container/operation/ContainerAttach
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
	return &LogFrame{
		Data: string(logBytes),
		Type: logType,
	}
}
