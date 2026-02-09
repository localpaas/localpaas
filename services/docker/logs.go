package docker

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/applog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

type ScanningLogOptions struct {
	BatchRecvOptions  batchrecvchan.Options
	ParseLogHeader    bool
	ParseLogTimestamp bool
}

type ScanningLogOption func(*ScanningLogOptions)

func WithParseLogHeader(flag bool) ScanningLogOption {
	return func(o *ScanningLogOptions) {
		o.ParseLogHeader = flag
	}
}

func WithParseLogTimestamp(flag bool) ScanningLogOption {
	return func(o *ScanningLogOptions) {
		o.ParseLogTimestamp = flag
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
) (logChan <-chan []*applog.LogFrame, closeFunc func() error) {
	opts := &ScanningLogOptions{
		ParseLogHeader: true,
	}
	for _, o := range options {
		o(opts)
	}

	batchChan := batchrecvchan.NewChan[*applog.LogFrame](opts.BatchRecvOptions)

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
			// Docker logs have an 8-byte header for stdout/stderr
			// Ref: https://docs.docker.com/reference/api/engine/version/v1.51/#tag/Container/operation/ContainerAttach
			_ = parseLogs(ctx, reader, batchChan, opts.ParseLogTimestamp)
		} else {
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				logFrame := &applog.LogFrame{
					Type: applog.LogTypeOut,
					Data: reflectutil.UnsafeBytesToStr(scanner.Bytes()),
				}
				if opts.ParseLogTimestamp {
					logFrame.ParseTimestampFromData()
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

// NOTE: This code is copied from `github.com/docker/docker/pkg/stdcopy`
// as we need to add some customization.

const (
	// Stdin represents standard input stream type.
	Stdin byte = iota
	// Stdout represents standard output stream type.
	Stdout
	// Stderr represents standard error steam type.
	Stderr
	// Systemerr represents errors originating from the system that make it
	// into the multiplexed stream.
	Systemerr

	stdWriterPrefixLen = 8
	stdWriterFdIndex   = 0
	stdWriterSizeIndex = 4

	startingBufLen = 32*1024 + stdWriterPrefixLen + 1
)

//nolint:gocognit
func parseLogs(
	ctx context.Context,
	src io.Reader,
	dst *batchrecvchan.Chan[*applog.LogFrame],
	parseTimestamp bool,
) (err error) {
	var (
		buf       = make([]byte, startingBufLen)
		bufLen    = len(buf)
		nr        int
		frameSize int
	)

	for {
		// Make sure we have at least a full header
		for nr < stdWriterPrefixLen {
			var nr2 int
			nr2, err = src.Read(buf[nr:])
			nr += nr2
			if errors.Is(err, io.EOF) {
				if nr < stdWriterPrefixLen {
					return nil
				}
				break
			}
			if err != nil {
				return apperrors.Wrap(err)
			}
		}

		logFrame := &applog.LogFrame{}

		// Check the first byte to know where to write
		switch buf[stdWriterFdIndex] {
		case Stdin:
			logFrame.Type = applog.LogTypeIn
		case Stdout:
			logFrame.Type = applog.LogTypeOut
		case Stderr, Systemerr:
			logFrame.Type = applog.LogTypeErr
		default:
			return fmt.Errorf("%w: Unrecognized input header: %d",
				apperrors.ErrInfraInternal, buf[stdWriterFdIndex])
		}

		// Retrieve the size of the frame
		frameSize = int(binary.BigEndian.Uint32(buf[stdWriterSizeIndex : stdWriterSizeIndex+4]))

		// Check if the buffer is big enough to read the frame.
		// Extend it if necessary.
		if frameSize+stdWriterPrefixLen > bufLen {
			buf = append(buf, make([]byte, frameSize+stdWriterPrefixLen-bufLen+1)...)
			bufLen = len(buf)
		}

		// While the amount of bytes read is less than the size of the frame + header, we keep reading
		for nr < frameSize+stdWriterPrefixLen {
			var nr2 int
			nr2, err = src.Read(buf[nr:])
			nr += nr2
			if errors.Is(err, io.EOF) {
				if nr < frameSize+stdWriterPrefixLen {
					return nil
				}
				break
			}
			if err != nil {
				return apperrors.Wrap(err)
			}
		}

		logFrame.Data = string(buf[stdWriterPrefixLen : frameSize+stdWriterPrefixLen])
		if parseTimestamp {
			logFrame.ParseTimestampFromData()
		}

		// Move the rest of the buffer to the beginning
		copy(buf, buf[frameSize+stdWriterPrefixLen:])
		// Move the index
		nr -= frameSize + stdWriterPrefixLen

		select {
		case <-ctx.Done(): // Make sure to quit if the context is done
			return nil
		default:
			dst.Send(logFrame)
		}
	}
}
