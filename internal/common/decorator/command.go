package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CommandHandler[C, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
	return QueryLoggingDecorator[C, R]{
		logger: logger,
		base: QueryMetricsDecorator[C, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
