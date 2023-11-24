package phylax

import (
	govsvc "github.com/deltaswapio/deltaswap-explorer/api/handlers/governor"
	heartbeatssvc "github.com/deltaswapio/deltaswap-explorer/api/handlers/heartbeats"
	vaasvc "github.com/deltaswapio/deltaswap-explorer/api/handlers/vaa"
	"github.com/deltaswapio/deltaswap-explorer/api/internal/config"
	"github.com/deltaswapio/deltaswap-explorer/api/routes/phylax/governor"
	"github.com/deltaswapio/deltaswap-explorer/api/routes/phylax/heartbeats"
	"github.com/deltaswapio/deltaswap-explorer/api/routes/phylax/phylax"
	"github.com/deltaswapio/deltaswap-explorer/api/routes/phylax/vaa"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RegisterRoutes sets up the handlers for the Phylax API.
func RegisterRoutes(
	cfg *config.AppConfig,
	app *fiber.App,
	rootLogger *zap.Logger,
	vaaService *vaasvc.Service,
	governorService *govsvc.Service,
	heartbeatsService *heartbeatssvc.Service,
) {

	// Set up controllers
	vaaCtrl := vaa.NewController(vaaService, rootLogger)
	governorCtrl := governor.NewController(governorService, rootLogger)
	phylaxCtrl := phylax.NewController(rootLogger, cfg.P2pNetwork)
	heartbeatsCtrl := heartbeats.NewController(heartbeatsService, rootLogger, cfg.P2pNetwork)

	// Set up route handlers
	apiV1 := app.Group("/v1")

	// signedVAA resource
	signedVAA := apiV1.Group("/signed_vaa")
	signedVAA.Get("/:chain/:emitter/:sequence", vaaCtrl.FindSignedVAAByID)
	signedBatchVAA := apiV1.Group("/signed_batch_vaa")
	signedBatchVAA.Get("/:chain/:trxID/:nonce", vaaCtrl.FindSignedBatchVAAByID)

	// phylaxSet resource
	phylaxSet := apiV1.Group("/phylaxset")
	phylaxSet.Get("/current", phylaxCtrl.GetPhylaxSet)

	// heartbeats resource
	heartbeats := apiV1.Group("/heartbeats")
	heartbeats.Get("", heartbeatsCtrl.GetLastHeartbeats)

	// governor resource
	gov := apiV1.Group("/governor")
	gov.Get("/available_notional_by_chain", governorCtrl.GetAvailNotionByChain)
	gov.Get("/enqueued_vaas", governorCtrl.GetEnqueuedVaas)
	gov.Get("/is_vaa_enqueued/:chain/:emitter/:sequence", governorCtrl.IsVaaEnqueued)
	gov.Get("/token_list", governorCtrl.GetTokenList)
}
