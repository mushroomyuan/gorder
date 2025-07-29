package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

// QueryHandler defines a generic type that recevices a Query Q,
// and returns a result R
type QueryHandler[Q, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
	return QueryLoggingDecorator[H, R]{
		logger: logger,
		base: QueryMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
