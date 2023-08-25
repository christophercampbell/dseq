package app

import (
	"context"

	"github.com/cometbft/cometbft/abci/types"
)

func (app *SequencerApplication) CheckTx(ctx context.Context, tx *types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	app.logger.Info("check tx", "app-id", app.ID)
	return &types.ResponseCheckTx{}, nil
}
