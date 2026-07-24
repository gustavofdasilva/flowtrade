package ports

import "context"

type TransactionManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
