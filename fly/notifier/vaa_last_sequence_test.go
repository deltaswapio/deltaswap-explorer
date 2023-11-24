package notifier

import (
	"os"
	"testing"

	"github.com/deltaswapio/deltaswap-explorer/fly/config"
	"github.com/test-go/testify/assert"
)

func TestNewLastSequenceNotifier(t *testing.T) {

	l := NewLastSequenceNotifier(nil, "mainnet-staging")

	assert.Equal(t, "mainnet-staging:wormscan:vaa-max-sequence", l.prefix)
}

func TestNewLastSequenceNotifierBackwardsCompat(t *testing.T) {

	prefix := config.GetPrefix()

	l := NewLastSequenceNotifier(nil, prefix)

	assert.Equal(t, "wormscan:vaa-max-sequence", l.prefix)
}

func TestNewLastSequenceNotifierWithPrefix(t *testing.T) {

	os.Setenv("P2P_NETWORK", "mainnet")
	os.Setenv("ENVIRONMENT", "staging")

	prefix := config.GetPrefix()

	l := NewLastSequenceNotifier(nil, prefix)

	assert.Equal(t, "mainnet-staging:wormscan:vaa-max-sequence", l.prefix)
}
