package base

import (
	"github.com/nats-io/go-nats-streaming"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConnect(t *testing.T) {
	_, err := Connect(stan.DefaultNatsURL, "stan-pub")
	require.NoError(t, err)
}
