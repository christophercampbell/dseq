package app

import (
	"context"
	"fmt"

	"github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/version"
)

func (app *SequencerApplication) Info(_ context.Context, _ *types.RequestInfo) (*types.ResponseInfo, error) {
	return &types.ResponseInfo{
		Data:             fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:          version.ABCIVersion,
		AppVersion:       AppVersion,
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.state.Hash(),
	}, nil
}
