package app

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/rand"
	"github.com/pkg/errors"
)

func (app *SequencerApplication) InitChain(_ context.Context, chain *types.RequestInitChain) (*types.ResponseInitChain, error) {
	app.logger.Info("initializing chain", "chain-id", chain.ChainId, "initial-height", chain.InitialHeight)
	return &types.ResponseInitChain{}, nil
}

func (app *SequencerApplication) PrepareProposal(_ context.Context, proposal *types.RequestPrepareProposal) (*types.ResponsePrepareProposal, error) {
	// simulate sequencing the transactions in some way...
	txs := make([][]byte, len(proposal.Txs))
	copy(txs, proposal.Txs)
	for i := len(txs) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		txs[i], txs[j] = txs[j], txs[i]
	}

	return &types.ResponsePrepareProposal{
		Txs: txs,
	}, nil
}

func (app *SequencerApplication) ProcessProposal(_ context.Context, proposal *types.RequestProcessProposal) (*types.ResponseProcessProposal, error) {
	return &types.ResponseProcessProposal{
		Status: types.ResponseProcessProposal_ACCEPT,
	}, nil
}

const (
	EtL2BlockStart datastreamer.EntryType = 1 // EtL2BlockStart entry type
	EtL2Tx         datastreamer.EntryType = 2 // EtL2Tx entry type
	EtL2BlockEnd   datastreamer.EntryType = 3 // EtL2BlockEnd entry type
	StSequencer                           = 1 // StSequencer sequencer stream type
)

func (app *SequencerApplication) FinalizeBlock(_ context.Context, block *types.RequestFinalizeBlock) (*types.ResponseFinalizeBlock, error) {
	app.stagedTxs = make([][]byte, 0)

	respTxs := make([]*types.ExecTxResult, len(block.Txs))

	if len(block.Txs) == 0 {
		return &types.ResponseFinalizeBlock{TxResults: respTxs, AppHash: app.state.Hash()}, nil
	}

	err := app.dataServer.StartAtomicOp()
	if err != nil {
		return nil, err
	}

	blockNum, err := app.dataServer.AddStreamEntry(EtL2BlockStart, []byte{})
	if err != nil {
		return nil, err
	}

	for i, tx := range block.Txs {

		app.stagedTxs = append(app.stagedTxs, tx)

		respTxs[i] = &types.ExecTxResult{
			Code: 0, // 0 == ok
			// TODO: potentially attach tx level events here as well
		}
		app.state.Size++

		_, err := app.dataServer.AddStreamEntry(EtL2Tx, tx)
		if err != nil {
			err = app.dataServer.RollbackAtomicOp()
			if err != nil {
				return nil, errors.Cause(err)
			}
			return nil, err
		}
	}

	_, err = app.dataServer.AddStreamEntry(EtL2BlockEnd, []byte{})
	if err != nil {
		app.logger.Error("error finalizing block to stream", "block", blockNum, "error", err)
		err = app.dataServer.RollbackAtomicOp()
		if err != nil {
			return nil, errors.Cause(err)
		}
		return nil, err
	}

	err = app.dataServer.CommitAtomicOp()
	if err != nil {
		return nil, err
	}

	app.state.Height = block.Height

	response := &types.ResponseFinalizeBlock{TxResults: respTxs, AppHash: app.state.Hash()} // hash should include tx hashes

	return response, nil
}

func (app *SequencerApplication) ExtendVote(_ context.Context, _ *types.RequestExtendVote) (*types.ResponseExtendVote, error) {
	return &types.ResponseExtendVote{}, nil
}

func (app *SequencerApplication) VerifyVoteExtension(_ context.Context, _ *types.RequestVerifyVoteExtension) (*types.ResponseVerifyVoteExtension, error) {
	return &types.ResponseVerifyVoteExtension{}, nil
}

func (app *SequencerApplication) Commit(_ context.Context, _ *types.RequestCommit) (*types.ResponseCommit, error) {
	if err := app.state.Save(); err != nil {
		app.logger.Error("app failed to save state", "error", err)
		return nil, err
	}

	return &types.ResponseCommit{}, nil
}
