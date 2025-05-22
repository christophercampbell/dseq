package app

import (
	"context"

	"github.com/cometbft/cometbft/abci/types"
)

func (app *SequencerApplication) Query(ctx context.Context, query *types.RequestQuery) (*types.ResponseQuery, error) {
	switch query.Path {
	default:
		return &types.ResponseQuery{}, nil
	}
}
