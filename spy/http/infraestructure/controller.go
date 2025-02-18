package infraestructure

import (
	"fmt"

	"github.com/deltaswapio/deltaswap-explorer/common/health"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Controller definition.
// Controller definition.
type Controller struct {
	checks []health.Check
	logger *zap.Logger
}

// NewController creates a Controller instance.
func NewController(checks []health.Check, logger *zap.Logger) *Controller {
	return &Controller{checks: checks, logger: logger}
}

// HealthCheck handler for the endpoint /health.
func (c *Controller) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(struct {
		Status string `json:"status"`
	}{Status: "OK"})
}

// ReadyCheck handler for the endpoint /ready.
func (c *Controller) ReadyCheck(ctx *fiber.Ctx) error {
	rctx := ctx.Context()
	requestID := fmt.Sprintf("%v", rctx.Value("requestid"))
	for _, check := range c.checks {
		if err := check(rctx); err != nil {
			c.logger.Error("Ready check failed", zap.Error(err), zap.String("requestID", requestID))
			return ctx.Status(fiber.StatusInternalServerError).JSON(struct {
				Ready string `json:"ready"`
				Error string `json:"error"`
			}{Ready: "NO", Error: err.Error()})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(struct {
		Ready string `json:"ready"`
	}{Ready: "OK"})

}
