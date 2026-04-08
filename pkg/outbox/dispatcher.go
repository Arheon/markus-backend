package outbox

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type Dispatcher struct {
	store  Store
	broker Broker

	interval        time.Duration
	readTimeout     time.Duration
	publishTimeout  time.Duration
	deleteTimeout   time.Duration
	updateTimeout   time.Duration
	maxMessages     int
	deleteBatchSize int
	maxAttempts     int32

	started         int32
	closed          int32
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	errCh           chan error
	discardedMsgsCh chan Message
	limit           int

	logger *logrus.Logger
}

type DispatcherOption func(*Dispatcher)

func WithInterval(interval time.Duration) DispatcherOption {
	return func(d *Dispatcher) {
		d.interval = interval
	}
}

func WithReadTimeout(readTimeout time.Duration) DispatcherOption {
	return func(d *Dispatcher) {
		d.readTimeout = readTimeout
	}
}

func WithPublishTimeout(publishTimeout time.Duration) DispatcherOption {
	return func(d *Dispatcher) {
		d.publishTimeout = publishTimeout
	}
}

func WithDeleteTimeout(deleteTimeout time.Duration) DispatcherOption {
	return func(d *Dispatcher) {
		d.deleteTimeout = deleteTimeout
	}
}

func WithUpdateTimeout(updateTimeout time.Duration) DispatcherOption {
	return func(d *Dispatcher) {
		d.updateTimeout = updateTimeout
	}
}

func WithReadBatchSize(maxMessages int) DispatcherOption {
	return func(d *Dispatcher) {
		if maxMessages > 0 {
			d.maxMessages = maxMessages
		}
	}
}

func WithDeleteBatchSize(deleteBatchSize int) DispatcherOption {
	return func(d *Dispatcher) {
		d.deleteBatchSize = deleteBatchSize
	}
}

func WithMaxAttempts(maxAttempts int32) DispatcherOption {
	return func(d *Dispatcher) {
		d.maxAttempts = maxAttempts
	}
}

func WithStarted(started int32) DispatcherOption {
	return func(d *Dispatcher) {
		d.started = started
	}
}

func WithClosed(closed int32) DispatcherOption {
	return func(d *Dispatcher) {
		d.closed = closed
	}
}

func WithCtx(ctx context.Context) DispatcherOption {
	return func(d *Dispatcher) {
		d.ctx = ctx
	}
}

func WithCancel(cancel context.CancelFunc) DispatcherOption {
	return func(d *Dispatcher) {
		d.cancel = cancel
	}
}

func WithErrChSize(size int) DispatcherOption {
	return func(d *Dispatcher) {
		d.errCh = make(chan error, size)
	}
}

func WithDiscardedMsgsChSize(size int) DispatcherOption {
	return func(d *Dispatcher) {
		d.discardedMsgsCh = make(chan Message, size)
	}
}

func WithLimit(limit int) DispatcherOption {
	return func(d *Dispatcher) {
		d.limit = limit
	}
}

func WithContext(ctx context.Context) DispatcherOption {
	return func(d *Dispatcher) {
		d.ctx = ctx
	}
}

func WithLogger(logger *logrus.Logger) DispatcherOption {
	return func(d *Dispatcher) {
		d.logger = logger
	}
}

func NewDispatcher(store Store, broker Broker, opts ...DispatcherOption) *Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())

	d := &Dispatcher{
		store:           store,
		broker:          broker,
		ctx:             ctx,
		cancel:          cancel,
		interval:        10 * time.Second,
		readTimeout:     5 * time.Second,
		publishTimeout:  5 * time.Second,
		deleteTimeout:   5 * time.Second,
		updateTimeout:   5 * time.Second,
		maxMessages:     100,
		deleteBatchSize: 20,
		maxAttempts:     math.MaxInt32,
	}

	for _, opt := range opts {
		opt(d)
	}

	if d.errCh == nil {
		d.errCh = make(chan error, 128)
	}

	if d.discardedMsgsCh == nil {
		d.discardedMsgsCh = make(chan Message, 128)
	}

	return d
}

func (d *Dispatcher) Start() {
	if d.logger != nil {
		d.logger.Debug("Start listening...")
	}

	if !atomic.CompareAndSwapInt32(&d.started, 0, 1) {
		return
	}

	d.wg.Add(1)
	go func() {
		d.logger.Debugf("Start ticker on %d seconds", d.interval)
		ticker := time.NewTicker(d.interval)

		defer d.wg.Done()
		defer close(d.errCh)
		defer close(d.discardedMsgsCh)
		defer ticker.Stop()

		for {
			select {
			case <-d.ctx.Done():
				return
			case <-ticker.C:
				d.logger.Debug("Process...")
				d.process()
				d.logger.Debug("Stop current process")
			}
		}
	}()
}

func (d *Dispatcher) process() {
	d.logger.Debug("Try fetch pending")
	msgs, err := d.store.FetchPending(d.ctx, d.limit)
	d.logger.WithFields(logrus.Fields{
		"messages": msgs,
		"error":    err,
	}).Debug("Messages received")
	if err != nil {
		d.logger.Errorf("receiving error while wetching messages... %v", err)
		d.errCh <- err
		return
	}

	for _, msg := range msgs {
		err := d.broker.Publish(d.ctx, msg)
		if err != nil {
			d.logger.Debugf("Mark message as failed... %v", err)
			_ = d.store.MarkFailed(d.ctx, msg.ID, err)
			select {
			case d.errCh <- err:
			case <-d.ctx.Done():
			default:
				d.logger.Error("Error channel is full, skipping...")
			}
			continue
		}

		d.logger.Debug("Process message...")
		_ = d.store.MarkProcessed(d.ctx, msg.ID)
		d.logger.Debug("Message processing success")
	}
}

func (d *Dispatcher) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&d.closed, 0, 1) {
		return nil
	}

	d.cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		d.wg.Wait()
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (d *Dispatcher) GetDiscardedMesgChan() chan Message {
	return d.discardedMsgsCh
}

func (d *Dispatcher) GetErrorChan() chan error {
	return d.errCh
}
