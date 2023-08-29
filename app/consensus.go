package app

import (
	"context"

	"github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/rand"
	"github.com/ethereum/go-ethereum/common"
)

func (app *SequencerApplication) InitChain(_ context.Context, chain *types.RequestInitChain) (*types.ResponseInitChain, error) {
	app.logger.Info("initializing chain", "chain-id", chain.ChainId, "initial-height", chain.InitialHeight)
	return &types.ResponseInitChain{}, nil
}

func (app *SequencerApplication) PrepareProposal(_ context.Context, proposal *types.RequestPrepareProposal) (*types.ResponsePrepareProposal, error) {
	app.logger.Debug("preparing proposal", "txs", len(proposal.Txs))

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
	proposer := common.BytesToAddress(proposal.ProposerAddress)
	app.logger.Debug("processing proposal", "proposer", proposer, "txs", len(proposal.Txs))
	return &types.ResponseProcessProposal{
		Status: types.ResponseProcessProposal_ACCEPT,
	}, nil
}

func (app *SequencerApplication) FinalizeBlock(_ context.Context, block *types.RequestFinalizeBlock) (*types.ResponseFinalizeBlock, error) {
	app.logger.Debug("finalize block", "height", block.Height, "txs", len(block.Txs), "hash", common.BytesToHash(block.Hash).Hex(), "size", block.Size())

	app.stagedTxs = make([][]byte, 0)

	// TODO: decide on punishment for block.Misbehavior[]

	// TODO: Execute transactions, accumulating []types.ExecTxResult

	respTxs := make([]*types.ExecTxResult, len(block.Txs))
	for i, tx := range block.Txs {

		app.stagedTxs = append(app.stagedTxs, tx)

		respTxs[i] = &types.ExecTxResult{
			Code: 0, // 0 == ok
			// TODO: potentially attach tx level events here as well
		}
		app.state.Size++

	}
	app.state.Height = block.Height

	response := &types.ResponseFinalizeBlock{TxResults: respTxs, AppHash: app.state.Hash()}

	// TODO: potentially raise block level events in response as well

	return response, nil
}

func (app *SequencerApplication) ExtendVote(_ context.Context, _ *types.RequestExtendVote) (*types.ResponseExtendVote, error) {
	app.logger.Info("extend vote")
	return &types.ResponseExtendVote{}, nil
}

func (app *SequencerApplication) VerifyVoteExtension(_ context.Context, _ *types.RequestVerifyVoteExtension) (*types.ResponseVerifyVoteExtension, error) {
	app.logger.Info("verify vote extension")
	return &types.ResponseVerifyVoteExtension{}, nil
}

func (app *SequencerApplication) Commit(_ context.Context, _ *types.RequestCommit) (*types.ResponseCommit, error) {
	app.logger.Info("commit")

	// TODO: apply the validator updates to state, this will require storing validator state, and tracking them in memory

	// FAKE: For the purposes of this POC, write the transactions in order to a file. Compare the sequences on different
	// nodes as proof that they come to sequence consensus
	for _, tx := range app.stagedTxs {
		app.sequence.Write(tx)
	}

	if err := app.state.Save(); err != nil {
		app.logger.Error("app failed to save state: %v", err)
		return nil, err
	}

	return &types.ResponseCommit{}, nil
}
