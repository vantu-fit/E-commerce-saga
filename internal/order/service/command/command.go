package command

import "context"

type CommandHanlder[C any] interface {
	Handle(ctx context.Context, cmd C) error
}
