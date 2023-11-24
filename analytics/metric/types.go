package metric

import (
	"context"

	"github.com/deltaswapio/deltaswap/sdk/vaa"
)

// MetricPushFunc is a function to push metrics
type MetricPushFunc func(context.Context, *vaa.VAA) error
