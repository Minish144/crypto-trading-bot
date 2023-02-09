package strategies

import "context"

type Strategy interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
