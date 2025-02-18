package alert

import (
	"fmt"

	"github.com/deltaswapio/deltaswap-explorer/common/client/alert"
)

// alert key constants definition.
const (
	ErrorSaveDestinationTx = "ERROR_SAVE_DESTINATION_TX"
)

func LoadAlerts(cfg alert.AlertConfig) map[string]alert.Alert {
	alerts := make(map[string]alert.Alert)

	// Alert error saving vaa.
	alerts[ErrorSaveDestinationTx] = alert.Alert{
		Alias:       ErrorSaveDestinationTx,
		Message:     fmt.Sprintf("[%s] %s", cfg.Environment, "Error saving destination tx in globalTransactions collection"),
		Description: "An error was found persisting the destination tx in globalTransactions collection.",
		Actions:     []string{"check globalTransactions collection"},
		Tags:        []string{cfg.Environment, "contract-watcher", "destination tx", "mongo"},
		Entity:      "contract-watcher",
		Priority:    alert.CRITICAL,
	}

	return alerts
}
