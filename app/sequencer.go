package app

import (
	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-hclog"
)

const (
	AppVersion uint64 = 1
)

type SequencerApplication struct {
	ID        string
	logger    hclog.Logger
	addr      common.Address
	state     *State
	sequence  *Sequence
	stagedTxs [][]byte

	// TODO: Store and maintain validator info for helping restarts, and punishing misbehavior
	// valAddrToPubKeyMap map[string]crypto.PublicKey
	// valUpdates []types.ValidatorUpdate
	dataServer *datastreamer.StreamServer
}

var _ types.Application = (*SequencerApplication)(nil)

func NewSequencer(logger hclog.Logger, identity string, addr common.Address, state *State, sequence *Sequence, ds *datastreamer.StreamServer) *SequencerApplication {
	return &SequencerApplication{
		ID:         identity,
		logger:     logger,
		addr:       addr,
		state:      state,
		sequence:   sequence,
		dataServer: ds,
	}
}
