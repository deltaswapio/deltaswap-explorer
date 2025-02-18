package vaa

import (
	"github.com/deltaswapio/deltaswap-explorer/analytics/metric"
	sdk "github.com/deltaswapio/deltaswap/sdk/vaa"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Controller controller struct definition.
type Controller struct {
	pushMetric metric.MetricPushFunc
	repository *Repository
	logger     *zap.Logger
}

// NewController create a new controller.
func NewController(pushMetric metric.MetricPushFunc, repository *Repository, logger *zap.Logger) *Controller {
	return &Controller{pushMetric: pushMetric, repository: repository, logger: logger}
}

// PushVAAMetrics push vaa metrics.
func (c *Controller) PushVAAMetrics(ctx *fiber.Ctx) error {
	payload := struct {
		ID string `json:"id"`
	}{}

	if err := ctx.BodyParser(&payload); err != nil {
		c.logger.Error("Error parsing request body", zap.Error(err))
		return err
	}

	c.logger.Info("Push VAA from endpoint", zap.String("id", payload.ID))

	vaaDoc, err := c.repository.FindById(ctx.Context(), payload.ID)
	if err != nil {
		c.logger.Error("Error finding VAA", zap.Error(err))
		return err
	}

	vaa, err := sdk.Unmarshal(vaaDoc.Vaa)
	if err != nil {
		c.logger.Error("Error unmarshalling VAA", zap.Error(err))
		return err
	}

	err = c.pushMetric(ctx.Context(), vaa)
	if err != nil {
		c.logger.Error("Error pushing metric", zap.Error(err))
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(struct {
		Push bool `json:"push"`
	}{Push: true})
}
