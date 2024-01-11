package app

import (
	"context"

	"github.com/cometbft/cometbft/abci/types"
)

func (app *SequencerApplication) CheckTx(ctx context.Context, tx *types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	return &types.ResponseCheckTx{Code: types.CodeTypeOK}, nil
}
