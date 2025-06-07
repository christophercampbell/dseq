package app

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/ethereum/go-ethereum/common"
)

const (
	AppVersion uint64 = 1
)

// SequencerApplication implements the ABCI application interface for the sequencer.
type SequencerApplication struct {
	types.BaseApplication

	ID        string
	logger    log.Logger
	addr      common.Address
	state     *State
	stagedTxs [][]byte

	// TODO: Store and maintain validator info for helping restarts, and punishing misbehavior
	// valAddrToPubKeyMap map[string]crypto.PublicKey
	// valUpdates []types.ValidatorUpdate
	dataServer *datastreamer.StreamServer
}

// Option configures a SequencerApplication.
type Option func(*SequencerApplication) error

// WithIdentity sets the application identity.
func WithIdentity(identity string) Option {
	return func(app *SequencerApplication) error {
		if identity == "" {
			return fmt.Errorf("identity cannot be empty")
		}
		app.ID = identity
		return nil
	}
}

// WithAddress sets the sequencer's address.
func WithAddress(addr common.Address) Option {
	return func(app *SequencerApplication) error {
		if addr == (common.Address{}) {
			return fmt.Errorf("address cannot be zero")
		}
		app.addr = addr
		return nil
	}
}

// WithState sets the application state.
func WithState(state *State) Option {
	return func(app *SequencerApplication) error {
		if state == nil {
			return fmt.Errorf("state cannot be nil")
		}
		app.state = state
		return nil
	}
}

// WithDataServer sets the data stream server.
func WithDataServer(ds *datastreamer.StreamServer) Option {
	return func(app *SequencerApplication) error {
		if ds == nil {
			return fmt.Errorf("data server cannot be nil")
		}
		app.dataServer = ds
		return nil
	}
}

var _ types.Application = (*SequencerApplication)(nil)

// NewSequencer constructs a SequencerApplication with the given logger and options.
// Returns an error if any of the options fail to apply.
func NewSequencer(logger log.Logger, opts ...Option) (*SequencerApplication, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	app := &SequencerApplication{
		logger: logger,
	}

	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return app, nil
}
