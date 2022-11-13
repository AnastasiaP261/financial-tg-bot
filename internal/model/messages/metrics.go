package messages

import (
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/financial-tg-bot/metrics"
)

var metricsCommAddPurchase = "add_purchase"

func metricsWrapper(wrappedFunc func() error, command string) error {
	if err := wrappedFunc(); err != nil {
		metrics.InFlightCommandMsg.WithLabelValues(command, metrics.StatusErr).Inc()
		return err
	}
	metrics.InFlightCommandMsg.WithLabelValues(command, metrics.StatusOk).Inc()
	return nil
}
