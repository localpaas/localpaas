package docker

import (
	"bufio"
	"context"
	"io"

	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/realtimelog"
)

func StartScanningLog(
	ctx context.Context,
	reader io.ReadCloser,
	options batchrecvchan.Options, // if zero, scan one by one
) (logChan <-chan []*realtimelog.LogFrame, closeFunc func() error) {
	batchChan := batchrecvchan.NewChan[*realtimelog.LogFrame](options)

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

func parseLogFrame(logBytes []byte) *realtimelog.LogFrame {
	var logType realtimelog.LogType
	// Format structure of the logs data, see:
	// https://docs.docker.com/reference/api/engine/version/v1.51/#tag/Container/operation/ContainerAttach
	//nolint:mnd
	if len(logBytes) > 8 {
		switch logBytes[0] {
		case 0:
			logType = realtimelog.LogTypeIn
		case 1:
			logType = realtimelog.LogTypeOut
		case 2:
			logType = realtimelog.LogTypeErr
		}
		logBytes = logBytes[8:]
	}
	return &realtimelog.LogFrame{
		Data: string(logBytes),
		Type: logType,
	}
}
