package phylax

import (
	"github.com/deltaswapio/deltaswap-explorer/api/handlers/phylax"
	"github.com/deltaswapio/deltaswap-explorer/api/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Controller definition.
type Controller struct {
	gs     phylax.PhylaxSet
	logger *zap.Logger
}

// NewController create a new controler.
func NewController(logger *zap.Logger, p2pNetwork string) *Controller {
	return &Controller{gs: phylax.GetByEnv(p2pNetwork),
		logger: logger.With(zap.String("module", "PhylaxController"))}
}

// PhylaxSetResponse response definition.
type PhylaxSetResponse struct {
	PhylaxSet PhylaxSet `json:"phylaxSet"`
}

// PhylaxSet response definition.
type PhylaxSet struct {
	Index     uint32   `json:"index"`
	Addresses []string `json:"addresses"`
}

// GetPhylaxSet godoc
// @Description Get current phylax set.
// @Tags Phylax
// @ID phylax-set
// @Success 200 {object} PhylaxSetResponse
// @Failure 400
// @Failure 500
// @Router /v1/phylaxset/current [get]
func (c *Controller) GetPhylaxSet(ctx *fiber.Ctx) error {
	// check phylaxSet exists.
	if len(c.gs.GstByIndex) == 0 {
		return response.NewApiError(ctx, fiber.StatusServiceUnavailable, response.Unavailable,
			"phylax set not fetched from chain yet", nil)
	}

	// get lasted phylaxSet.
	guardinSet := c.gs.GetLatest()

	// get phylax addresses.
	addresses := make([]string, len(guardinSet.Keys))
	for i, v := range guardinSet.Keys {
		addresses[i] = v.Hex()
	}

	// create response.
	response := PhylaxSetResponse{
		PhylaxSet: PhylaxSet{
			Index:     guardinSet.Index,
			Addresses: addresses,
		},
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}
