package app

import (
	"net"
	"os"
	"testing"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getFreePort finds an available TCP port and returns it.
func getFreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func setupTestSequencer(t *testing.T) (*SequencerApplication, func()) {
	// Create temporary directory for test data
	tmpDir, err := os.MkdirTemp("", "sequencer_test")
	require.NoError(t, err)

	// Create state
	state, err := NewState(tmpDir)
	require.NoError(t, err)

	// Get a free port for the data server
	port := getFreePort(t)

	// Create data server
	ds, err := datastreamer.NewServer(
		uint16(port),               // Use dynamic port
		0,                          // fileFormat
		0,                          // maxFileSize
		datastreamer.StreamType(1), // streamType
		tmpDir,                     // path
		nil,                        // logConfig
	)
	require.NoError(t, err)

	// Start the server
	err = ds.Start()
	require.NoError(t, err)

	// Create logger
	logger := log.NewNopLogger()

	// Create sequencer
	app, err := NewSequencer(logger,
		WithIdentity("test-sequencer"),
		WithAddress(common.HexToAddress("0x1234")),
		WithState(state),
		WithDataServer(ds),
	)
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		app.state.Close()
		os.RemoveAll(tmpDir)
	}

	return app, cleanup
}

func TestNewSequencer(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name: "valid options",
			opts: []Option{
				WithIdentity("test"),
				WithAddress(common.HexToAddress("0x1234")),
			},
			wantErr: false,
		},
		{
			name: "empty identity",
			opts: []Option{
				WithIdentity(""),
				WithAddress(common.HexToAddress("0x1234")),
			},
			wantErr: true,
		},
		{
			name: "zero address",
			opts: []Option{
				WithIdentity("test"),
				WithAddress(common.Address{}),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := NewSequencer(log.NewNopLogger(), tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, app)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, app)
			}
		})
	}
}

func TestSequencerOptions(t *testing.T) {
	app, cleanup := setupTestSequencer(t)
	defer cleanup()

	// Test identity
	assert.Equal(t, "test-sequencer", app.ID)

	// Test address
	assert.Equal(t, common.HexToAddress("0x1234"), app.addr)

	// Test state
	assert.NotNil(t, app.state)
	assert.Equal(t, int64(0), app.state.Size)
	assert.Equal(t, int64(0), app.state.Height)

	// Test data server
	assert.NotNil(t, app.dataServer)
}

func TestSequencerStagedTransactions(t *testing.T) {
	app, cleanup := setupTestSequencer(t)
	defer cleanup()

	// Test initial state
	assert.Empty(t, app.stagedTxs)

	// Add staged transactions
	tx1 := []byte("tx1")
	tx2 := []byte("tx2")
	app.stagedTxs = append(app.stagedTxs, tx1, tx2)

	// Verify staged transactions
	assert.Len(t, app.stagedTxs, 2)
	assert.Equal(t, tx1, app.stagedTxs[0])
	assert.Equal(t, tx2, app.stagedTxs[1])
}

func TestSequencerWithNilLogger(t *testing.T) {
	app, err := NewSequencer(nil)
	assert.Error(t, err)
	assert.Nil(t, app)
}

func TestSequencerWithNilState(t *testing.T) {
	app, err := NewSequencer(log.NewNopLogger(), WithState(nil))
	assert.Error(t, err)
	assert.Nil(t, app)
}

func TestSequencerWithNilDataServer(t *testing.T) {
	app, err := NewSequencer(log.NewNopLogger(), WithDataServer(nil))
	assert.Error(t, err)
	assert.Nil(t, app)
}
