package core

import "context"

type Stream interface {
	Send(data any) error
	Close() error
	Context() context.Context
}

type StreamHandler func(ctx context.Context, req any, stream Stream) error

type StreamHandlerFunc[Req any] func(ctx context.Context, req *Req, stream Stream) error
