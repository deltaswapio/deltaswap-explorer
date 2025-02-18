package consumer

import (
	"context"

	"github.com/deltaswapio/deltaswap-explorer/parser/internal/metrics"
	"github.com/deltaswapio/deltaswap-explorer/parser/processor"
	"github.com/deltaswapio/deltaswap-explorer/parser/queue"
	"go.uber.org/zap"
)

// Consumer consumer struct definition.
type Consumer struct {
	consume queue.VAAConsumeFunc
	process processor.ProcessorFunc
	metrics metrics.Metrics
	logger  *zap.Logger
}

// New creates a new vaa consumer.
func New(consume queue.VAAConsumeFunc, process processor.ProcessorFunc, metrics metrics.Metrics, logger *zap.Logger) *Consumer {
	return &Consumer{consume: consume, process: process, metrics: metrics, logger: logger}
}

// Start consumes messages from VAA queue, parse and store those messages in a repository.
func (c *Consumer) Start(ctx context.Context) {
	go func() {
		for msg := range c.consume(ctx) {
			event := msg.Data()

			// check id message is expired.
			if msg.IsExpired() {
				c.logger.Warn("Message with vaa expired", zap.String("id", event.ID))
				msg.Failed()
				continue
			}
			c.metrics.IncVaaUnexpired(event.ChainID)

			_, err := c.process(ctx, event.Vaa)
			if err != nil {
				c.logger.Error("Error processing parsed vaa",
					zap.String("id", event.ID),
					zap.Error(err))
				msg.Failed()
				continue
			}
			msg.Done()
		}
	}()
}
