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

var _ types.Application = (*SequencerApplication)(nil)

func NewSequencer(logger log.Logger, identity string, addr common.Address, state *State, ds *datastreamer.StreamServer) *SequencerApplication {
	return &SequencerApplication{
		ID:         identity,
		logger:     logger,
		addr:       addr,
		state:      state,
		dataServer: ds,
	}
}
