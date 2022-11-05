package types

import "context"

type BaseService interface {
	Start(context.Context) error
	Stop() error
}

type Service interface {
	BaseService
	Name() string
	Healthy() bool
	ExtWrite() chan<- Message // channel for passing messages to external services
}

type Message struct {
	From string
	To   string
	Data interface{}
}

// OnTick - handler run when ticks in the base service occur
// generates a message to send
type OnTick func() Message

// ProcessMsg - handler run when a message in the base service is received
type ProcessMsg func(Message)