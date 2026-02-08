package docker

import (
	"bufio"
	"context"
	"io"

	"github.com/docker/docker/pkg/stdcopy"

	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/realtimelog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

type ScanningLogOptions struct {
	BatchRecvOptions batchrecvchan.Options
	ParseLogHeader   bool
}

type ScanningLogOption func(*ScanningLogOptions)

func WithParseLogHeader(flag bool) ScanningLogOption {
	return func(o *ScanningLogOptions) {
		o.ParseLogHeader = flag
	}
}

func WithBatchRecvOptions(recvOpts batchrecvchan.Options) ScanningLogOption {
	return func(o *ScanningLogOptions) {
		o.BatchRecvOptions = recvOpts
	}
}

type logWriter struct {
	LogType realtimelog.LogType
	LogChan *batchrecvchan.Chan[*realtimelog.LogFrame]
}

func (w *logWriter) Write(p []byte) (int, error) {
	w.LogChan.Send(&realtimelog.LogFrame{
		Type: w.LogType,
		Data: string(p),
	})
	return len(p), nil
}

func StartScanningLog(
	ctx context.Context,
	reader io.ReadCloser,
	options ...ScanningLogOption,
) (logChan <-chan []*realtimelog.LogFrame, closeFunc func() error) {
	opts := &ScanningLogOptions{
		ParseLogHeader: true,
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

		if opts.ParseLogHeader {
			outWriter := &logWriter{
				LogType: realtimelog.LogTypeOut,
				LogChan: batchChan,
			}
			errWriter := &logWriter{
				LogType: realtimelog.LogTypeErr,
				LogChan: batchChan,
			}
			// Use stdcopy.StdCopy to demultiplex and format the logs
			// Docker logs have an 8-byte header for stdout/stderr which stdcopy handles
			// Ref: https://docs.docker.com/reference/api/engine/version/v1.51/#tag/Container/operation/ContainerAttach
			_, err := stdcopy.StdCopy(outWriter, errWriter, reader)
			if err != nil {
				panic(err)
			}
		} else {
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				logFrame := &realtimelog.LogFrame{
					Type: realtimelog.LogTypeOut,
					Data: reflectutil.UnsafeBytesToStr(scanner.Bytes()),
				}
				select {
				case <-ctx.Done(): // Make sure to quit if the context is done
					return
				default:
					batchChan.Send(logFrame)
				}
			}
		}
	}()

	return batchChan.Receiver(), func() error { return reader.Close() }
}
