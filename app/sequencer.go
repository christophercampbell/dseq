package app

import (
	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/ethereum/go-ethereum/common"
)

const (
	AppVersion uint64 = 1
)

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
type Option func(*SequencerApplication)

// WithIdentity sets the application identity.
func WithIdentity(identity string) Option {
	return func(app *SequencerApplication) {
		app.ID = identity
	}
}

// WithAddress sets the sequencer's address.
func WithAddress(addr common.Address) Option {
	return func(app *SequencerApplication) {
		app.addr = addr
	}
}

// WithState sets the application state.
func WithState(state *State) Option {
	return func(app *SequencerApplication) {
		app.state = state
	}
}

// WithDataServer sets the data stream server.
func WithDataServer(ds *datastreamer.StreamServer) Option {
	return func(app *SequencerApplication) {
		app.dataServer = ds
	}
}

var _ types.Application = (*SequencerApplication)(nil)

// NewSequencer constructs a SequencerApplication with the given logger and options.
func NewSequencer(logger log.Logger, opts ...Option) *SequencerApplication {
	app := &SequencerApplication{
		logger: logger,
	}
	for _, opt := range opts {
		opt(app)
	}
	return app
}
