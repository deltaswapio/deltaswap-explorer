package processor

import (
	"context"

	"github.com/deltaswapio/deltaswap-explorer/parser/parser"
)

// ProcessorFunc is a function to process vaa message.
type ProcessorFunc func(context.Context, []byte) (*parser.ParsedVaaUpdate, error)
