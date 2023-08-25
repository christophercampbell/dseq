package app

import (
	"context"

	"github.com/cometbft/cometbft/abci/types"
)

func (app *SequencerApplication) ListSnapshots(ctx context.Context, snapshots *types.RequestListSnapshots) (*types.ResponseListSnapshots, error) {
	//TODO implement me
	panic("implement me")
}

func (app *SequencerApplication) OfferSnapshot(ctx context.Context, snapshot *types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error) {
	//TODO implement me
	panic("implement me")
}

func (app *SequencerApplication) LoadSnapshotChunk(ctx context.Context, chunk *types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error) {
	//TODO implement me
	panic("implement me")
}

func (app *SequencerApplication) ApplySnapshotChunk(ctx context.Context, chunk *types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error) {
	//TODO implement me
	panic("implement me")
}
