package batchrecvchan

import (
	"sync"
	"time"

	"github.com/tiendc/gofn"
)

const (
	minThresholdPeriod = 50 * time.Millisecond
	defaultChanSize    = 100
)

// Options sends data when either the period of time passes or the item count
// reaches the max value.
type Options struct {
	ThresholdPeriod time.Duration
	MaxItem         int
	ChanSize        int // 100 if unset
}

type Chan[T any] struct {
	thresholdPeriod time.Duration
	maxItem         int

	channel      chan []T
	mu           sync.Mutex
	currentBatch []T
	batchMode    bool
	stopped      bool
}

func (ch *Chan[T]) Send(item T) {
	if !ch.batchMode {
		ch.channel <- []T{item}
		return
	}

	ch.mu.Lock()
	ch.currentBatch = append(ch.currentBatch, item)
	sendData := len(ch.currentBatch) >= ch.maxItem
	ch.mu.Unlock()
	if sendData {
		ch.sendData()
	}
}

func (ch *Chan[T]) sendData() {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	if len(ch.currentBatch) > 0 {
		ch.channel <- ch.currentBatch
		ch.currentBatch = ch.currentBatch[:0]
	}
}

func (ch *Chan[T]) Receiver() <-chan []T {
	return ch.channel
}

func (ch *Chan[T]) Close() {
	ch.stopped = true
	close(ch.channel)
}

func NewChan[T any](options Options) *Chan[T] {
	ch := &Chan[T]{
		thresholdPeriod: options.ThresholdPeriod,
		maxItem:         options.MaxItem,
	}
	if ch.thresholdPeriod > 0 && ch.thresholdPeriod < minThresholdPeriod {
		ch.thresholdPeriod = minThresholdPeriod
	}
	if ch.thresholdPeriod == 0 && ch.maxItem == 0 {
		ch.maxItem = 1
	}
	ch.batchMode = ch.thresholdPeriod > 0 || ch.maxItem > 1
	if ch.batchMode {
		ch.currentBatch = make([]T, 0, ch.maxItem)
	}
	ch.channel = make(chan []T, gofn.Coalesce(options.ChanSize, defaultChanSize))

	if ch.thresholdPeriod > 0 {
		go func() {
			for range time.Tick(ch.thresholdPeriod) {
				ch.sendData()
				if ch.stopped {
					return
				}
			}
		}()
	}
	return ch
}
