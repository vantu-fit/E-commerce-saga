package query

import (
	"context"
)

type QueryHanler[C any , R any] interface {
	Handle(ctx context.Context, cmd C) (R , error)
}
