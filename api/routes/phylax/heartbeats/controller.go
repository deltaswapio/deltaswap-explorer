package heartbeats

import (
	"strconv"

	"github.com/deltaswapio/deltaswap-explorer/api/handlers/heartbeats"
	"github.com/deltaswapio/deltaswap-explorer/api/handlers/phylax"
	"github.com/deltaswapio/deltaswap-explorer/api/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Controller definition.
type Controller struct {
	srv    *heartbeats.Service
	logger *zap.Logger
	gs     phylax.PhylaxSet
}

// NewController create a new controler.
func NewController(srv *heartbeats.Service, logger *zap.Logger, p2pNetwork string) *Controller {
	return &Controller{
		srv:    srv,
		logger: logger.With(zap.String("module", "HeartbeatsController")),
		gs:     phylax.GetByEnv(p2pNetwork),
	}
}

// HeartbeatsResponse response.
type HeartbeatsResponse struct {
	Heartbeats []*HeartbeatResponse `json:"entries"`
}

type HeartbeatResponse struct {
	VerifiedPhylaxAddr string        `json:"verifiedPhylaxAddr"`
	P2PNodeAddr        string        `json:"p2pNodeAddr"`
	RawHeartbeat       *RawHeartbeat `json:"rawHeartbeat"`
}

type RawHeartbeat struct {
	NodeName      string                      `json:"nodeName"`
	Counter       string                      `json:"counter"`
	Timestamp     string                      `json:"timestamp"`
	Networks      []*HeartbeatNetworkResponse `json:"networks"`
	Version       string                      `json:"version"`
	PhylaxAddr    string                      `json:"phylaxAddr"`
	BootTimestamp string                      `json:"bootTimestamp"`
	Features      []string                    `json:"features"`
}

// HeartbeatNetwork definition.
type HeartbeatNetworkResponse struct {
	ID              int64  `bson:"id" json:"id"`
	Height          string `bson:"height" json:"height"`
	ContractAddress string `bson:"contractaddress" json:"contractAddress"`
	ErrorCount      string `bson:"errorcount" json:"errorCount"`
}

// GetPhylaxSet godoc
// @Description Get heartbeats for phylaxs
// @Tags Phylax
// @ID phylaxs-hearbeats
// @Success 200 {object} HeartbeatsResponse
// @Failure 400
// @Failure 500
// @Router /v1/heartbeats [get]
func (c *Controller) GetLastHeartbeats(ctx *fiber.Ctx) error {

	// check phylaxSet exists.
	if len(c.gs.GstByIndex) == 0 {
		err := response.NewApiError(
			ctx,
			fiber.StatusServiceUnavailable,
			response.Unavailable,
			"phylax set not fetched from chain yet",
			nil,
		)
		return err
	}

	// get the latest phylaxSet.
	phylaxSet := c.gs.GetLatest()
	phylaxAddresses := phylaxSet.KeysAsHexStrings()

	// get last heartbeats by ids.
	heartbeats, err := c.srv.GetHeartbeatsByIds(ctx.Context(), phylaxAddresses)
	if err != nil {
		return err
	}

	// build heartbeats response compatible with grpc api response.
	response := buildHeartbeatResponse(heartbeats)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func buildHeartbeatResponse(heartbeats []*heartbeats.HeartbeatDoc) *HeartbeatsResponse {
	if heartbeats == nil {
		return nil
	}
	heartbeatResponses := make([]*HeartbeatResponse, 0, len(heartbeats))
	for _, heartbeat := range heartbeats {

		networkResponses := make([]*HeartbeatNetworkResponse, 0, len(heartbeat.Networks))
		for _, network := range heartbeat.Networks {
			networkResponse := &HeartbeatNetworkResponse{
				ID:              network.ID,
				Height:          strconv.Itoa(int(network.Height)),
				ContractAddress: network.ContractAddress,
				ErrorCount:      strconv.Itoa(int(network.ErrorCount)),
			}
			networkResponses = append(networkResponses, networkResponse)
		}

		hr := HeartbeatResponse{
			VerifiedPhylaxAddr: heartbeat.ID,
			P2PNodeAddr:        "", // not exists in heartbeats mongo collection.
			RawHeartbeat: &RawHeartbeat{
				NodeName:      heartbeat.NodeName,
				Counter:       strconv.Itoa(int(heartbeat.Counter)),
				Timestamp:     strconv.Itoa(int(heartbeat.Timestamp)),
				Networks:      networkResponses,
				Version:       heartbeat.Version,
				PhylaxAddr:    heartbeat.PhylaxAddr,
				BootTimestamp: strconv.Itoa(int(heartbeat.BootTimestamp)),
				Features:      heartbeat.Features,
			},
		}
		heartbeatResponses = append(heartbeatResponses, &hr)
	}
	return &HeartbeatsResponse{
		Heartbeats: heartbeatResponses,
	}
}
