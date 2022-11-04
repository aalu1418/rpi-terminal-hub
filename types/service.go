package types

import "context"

type Service interface {
	Name() string
	Start(context.Context) error
	Healthy() bool
	ExtRead() <-chan Message  // channel for passing external messages to service
	ExtWrite() chan<- Message // channel for passing messages to external services
	Stop() error
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