package query

import "context"

type QueryHandler[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}
