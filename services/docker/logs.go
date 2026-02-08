package docker

import (
	"bufio"
	"context"
	"io"

	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/realtimelog"
)

type ScanningLogOptions struct {
	BatchRecvOptions batchrecvchan.Options
	ParseFrameHeader bool
}

type ScanningLogOption func(*ScanningLogOptions)

func WithParseFrameHeader(flag bool) ScanningLogOption {
	return func(o *ScanningLogOptions) {
		o.ParseFrameHeader = flag
	}
}

func WithBatchRecvOptions(recvOpts batchrecvchan.Options) ScanningLogOption {
	return func(o *ScanningLogOptions) {
		o.BatchRecvOptions = recvOpts
	}
}

func StartScanningLog(
	ctx context.Context,
	reader io.ReadCloser,
	options ...ScanningLogOption,
) (logChan <-chan []*realtimelog.LogFrame, closeFunc func() error) {
	opts := &ScanningLogOptions{
		ParseFrameHeader: true,
	}
	for _, o := range options {
		o(opts)
	}

	batchChan := batchrecvchan.NewChan[*realtimelog.LogFrame](opts.BatchRecvOptions)

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
			logFrame := parseLogFrame(scanner.Bytes(), opts.ParseFrameHeader)
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

func parseLogFrame(logBytes []byte, parseHeader bool) *realtimelog.LogFrame {
	var logType realtimelog.LogType
	// Format structure of the logs data, see:
	// https://docs.docker.com/reference/api/engine/version/v1.51/#tag/Container/operation/ContainerAttach
	//nolint:mnd
	if parseHeader && len(logBytes) > 8 {
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
