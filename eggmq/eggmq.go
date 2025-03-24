package eggmq

import (
	"context"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/sys"
	"go.uber.org/zap"
	"sync"
	"time"
)

type EggMQ struct {
	code       string
	ch         chan *Message
	maxRetries int
	retryDelay time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.RWMutex
	listeners  map[string][]Listener
}

func NewEggMQ(
	code string,
	bufferSize int,
	maxRetries int,
	retryDelay time.Duration,
) *EggMQ {
	if bufferSize <= 0 {
		bufferSize = 1024
	}
	if maxRetries <= 0 {
		maxRetries = 3
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &EggMQ{
		code:       code,
		ch:         make(chan *Message, bufferSize),
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		ctx:        ctx,
		cancel:     cancel,
		listeners:  map[string][]Listener{},
	}
}

func (mq *EggMQ) Code() string {
	return mq.code
}

func (mq *EggMQ) Startup() *errors.Error {
	go mq.loadPendingMessages()
	go mq.doStartup()
	return nil
}

func (mq *EggMQ) Shutdown(ctx context.Context) *errors.Error {
	mq.cancel()
	mq.wg.Wait()
	close(mq.ch)
	return nil
}

func (mq *EggMQ) Send(id string, topic string, payload string) *errors.Error {
	msg := NewMessage(id, topic, payload)
	err := MessageCreate(msg)
	if err != nil {
		return err
	}
	mq.ch <- msg
	return nil
}

func (mq *EggMQ) RegisterListener(listener Listener) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	topic := listener.Topic()
	list, ok := mq.listeners[topic]
	if !ok {
		list = []Listener{}
	}
	list = append(list, listener)
	mq.listeners[topic] = list
	sys.Info("#", mq.Code(), "# register listener: #", topic, "#")
	gLogger.Info("Register Listener: ", zap.String("topic", topic))
}

func (mq *EggMQ) getListeners(topic string) []Listener {
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	list, ok := mq.listeners[topic]
	if !ok {
		return []Listener{}
	}
	return list
}

func (mq *EggMQ) doStartup() {
	for {
		select {
		case msg := <-mq.ch:
			mq.doDeal(msg)
		case <-mq.ctx.Done():
			return
		}
	}
}

func (mq *EggMQ) doDeal(msg *Message) *errors.Error {
	mq.mu.Lock()
	err := MessageToProcessing(msg)
	if err != nil {
		mq.mu.Unlock()
		return err
	}
	mq.mu.Unlock()

	err = mq.doDealMsg(msg)

	mq.mu.Lock()
	defer mq.mu.Unlock()

	if err != nil {
		if msg.RetryCount >= mq.maxRetries {
			MessageToFailed(msg)
			//TODO add retry for err
		} else {
			time.AfterFunc(mq.retryDelay, func() {
				mq.ch <- msg
			})
			MessageToPending(msg)
			//TODO add retry for err
		}
		return nil
	}

	MessageToCompleted(msg)
	//TODO add retry for err
	return nil
}

func (mq *EggMQ) doDealMsg(msg *Message) *errors.Error {
	listeners := mq.getListeners(msg.Topic)
	for _, listener := range listeners {
		err := listener.Handle(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mq *EggMQ) loadPendingMessages() {
	MessageLoadPending(func(msg *Message) {
		mq.ch <- msg
	})
	//TODO add retry for err
}
