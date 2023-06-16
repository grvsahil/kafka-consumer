package consumer

import "context"

type Consume interface {
	Consume(ctx context.Context, topics []string) error
}
