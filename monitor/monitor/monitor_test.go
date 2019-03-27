package monitor

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPortMonitor_Monitor(t *testing.T) {
	m:=NewPortMonitor("47.52.173.202","1234")
	require.NoError(t,m.Monitor())
}
