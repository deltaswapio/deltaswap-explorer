package deltaswapscan

import (
	"time"

	addrsvc "github.com/deltaswapio/deltaswap-explorer/api/handlers/address"
	govsvc "github.com/deltaswapio/deltaswap-explorer/api/handlers/governor"
	infrasvc "github.com/deltaswapio/deltaswap-explorer/api/handlers/infrastructure"
	obssvc "github.com/deltaswapio/deltaswap-explorer/api/handlers/observations"
	relayssvc "github.com/deltaswapio/deltaswap-explorer/api/handlers/relays"
	trxsvc "github.com/deltaswapio/deltaswap-explorer/api/handlers/transactions"
	vaasvc "github.com/deltaswapio/deltaswap-explorer/api/handlers/vaa"
	"github.com/deltaswapio/deltaswap-explorer/api/routes/deltaswapscan/address"
	"github.com/deltaswapio/deltaswap-explorer/api/routes/deltaswapscan/governor"
	"github.com/deltaswapio/deltaswap-explorer/api/routes/deltaswapscan/infrastructure"
	"github.com/deltaswapio/deltaswap-explorer/api/routes/deltaswapscan/observations"
	"github.com/deltaswapio/deltaswap-explorer/api/routes/deltaswapscan/relays"
	"github.com/deltaswapio/deltaswap-explorer/api/routes/deltaswapscan/transactions"
	"github.com/deltaswapio/deltaswap-explorer/api/routes/deltaswapscan/vaa"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"go.uber.org/zap"
)

var cacheConfig = cache.Config{
	Next: func(c *fiber.Ctx) bool {
		return c.Query("refresh") == "true"
	},
	Expiration:           1 * time.Second,
	CacheControl:         true,
	StoreResponseHeaders: true,
}

// RegisterRoutes sets up the handlers for the Wormscan API.
func RegisterRoutes(
	app *fiber.App,
	rootLogger *zap.Logger,
	addressService *addrsvc.Service,
	vaaService *vaasvc.Service,
	obsService *obssvc.Service,
	governorService *govsvc.Service,
	infrastructureService *infrasvc.Service,
	transactionsService *trxsvc.Service,
	relaysService *relayssvc.Service,
) {

	// Set up controllers
	addressCtrl := address.NewController(addressService, rootLogger)
	vaaCtrl := vaa.NewController(vaaService, rootLogger)
	observationsCtrl := observations.NewController(obsService, rootLogger)
	governorCtrl := governor.NewController(governorService, rootLogger)
	infrastructureCtrl := infrastructure.NewController(infrastructureService)
	transactionCtrl := transactions.NewController(transactionsService, rootLogger)
	relaysCtrl := relays.NewController(relaysService, rootLogger)

	// Set up route handlers
	api := app.Group("/api/v1")
	api.Use(cors.New()) // TODO CORS restrictions?

	// monitoring
	api.Get("/health", infrastructureCtrl.HealthCheck)
	api.Get("/ready", infrastructureCtrl.ReadyCheck)
	api.Get("/version", infrastructureCtrl.Version)

	// accounts resource
	api.Get("/address/:id", addressCtrl.FindById)

	// analytics, transactions, custom endpoints
	api.Get("/global-tx/:chain/:emitter/:sequence", transactionCtrl.FindGlobalTransactionByID)
	api.Get("/last-txs", transactionCtrl.GetLastTransactions)
	api.Get("/scorecards", transactionCtrl.GetScorecards)
	api.Get("/x-chain-activity", transactionCtrl.GetChainActivity)
	api.Get("/top-assets-by-volume", transactionCtrl.GetTopAssets)
	api.Get("/top-chain-pairs-by-num-transfers", transactionCtrl.GetTopChainPairs)
	api.Get("token/:chain/:token_address", transactionCtrl.GetTokenByChainAndAddress)
	api.Get("/transactions", transactionCtrl.ListTransactions)
	api.Get("/transactions/:chain/:emitter/:sequence", transactionCtrl.GetTransactionByID)

	// vaas resource
	vaas := api.Group("/vaas")
	vaas.Use(cache.New(cacheConfig))
	vaas.Get("/vaa-counts", vaaCtrl.GetVaaCount)
	vaas.Get("/", vaaCtrl.FindAll)
	vaas.Get("/:chain", vaaCtrl.FindByChain)
	vaas.Get("/:chain/:emitter", vaaCtrl.FindByEmitter)
	vaas.Get("/:chain/:emitter/:sequence", vaaCtrl.FindById)
	vaas.Post("/parse", vaaCtrl.ParseVaa)

	// oservations resource
	observations := api.Group("/observations")
	observations.Get("/", observationsCtrl.FindAll)
	observations.Get("/:chain", observationsCtrl.FindAllByChain)
	observations.Get("/:chain/:emitter", observationsCtrl.FindAllByEmitter)
	observations.Get("/:chain/:emitter/:sequence", observationsCtrl.FindAllByVAA)
	observations.Get("/:chain/:emitter/:sequence/:signer/:hash", observationsCtrl.FindOne)

	// governor resources
	governor := api.Group("/governor")
	governorLimit := governor.Group("/limit")
	governorLimit.Get("/", governorCtrl.GetGovernorLimit)

	governorConfigs := governor.Group("/config")
	governorConfigs.Get("/", governorCtrl.FindGovernorConfigurations)
	governorConfigs.Get("/:phylax_address", governorCtrl.FindGovernorConfigurationByPhylaxAddress)

	governorStatus := governor.Group("/status")
	governorStatus.Get("/", governorCtrl.FindGovernorStatus)
	governorStatus.Get("/:phylax_address", governorCtrl.FindGovernorStatusByPhylaxAddress)

	governorNotional := governor.Group("/notional")
	governorNotional.Get("/limit/", governorCtrl.FindNotionalLimit)
	governorNotional.Get("/limit/:chain", governorCtrl.GetNotionalLimitByChainID)
	governorNotional.Get("/available/", governorCtrl.GetAvailableNotional)
	governorNotional.Get("/available/:chain", governorCtrl.GetAvailableNotionalByChainID)
	governorNotional.Get("/max_available/:chain", governorCtrl.GetMaxNotionalAvailableByChainID)

	enqueueVaas := governor.Group("/enqueued_vaas")
	enqueueVaas.Get("/", governorCtrl.GetEnqueuedVaas)
	enqueueVaas.Get("/:chain", governorCtrl.GetEnqueuedVaasByChainID)

	relays := api.Group("/relays")
	relays.Get("/:chain/:emitter/:sequence", relaysCtrl.FindOne)
}
