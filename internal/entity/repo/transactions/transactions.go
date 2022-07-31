package transactions

import "context"

type TransactionInterface interface {
	Start(ctx context.Context) (context.Context, error)
	Finish(ctx context.Context, err *error)
}
