package server

import (
	"context"
	"go.uber.org/zap"
)

type Storager interface {
	UseLogger(logger *zap.Logger)

	InsertChan(ctx context.Context, in <-chan *Response)
}

type Requester interface {
	UseLogger(logger *zap.Logger)

	GetBatch(ctx context.Context, out chan<- *Response, batch RequestBatch) int
}

type StorageConfiger interface {
	Build() (Storager, error)
}

type RequesterConfiger interface {
	Build() (Requester, error)
}
