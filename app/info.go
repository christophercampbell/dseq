package app

import (
	"context"
	"encoding/json"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/version"
)

func (app *SequencerApplication) Info(_ context.Context, _ *types.RequestInfo) (*types.ResponseInfo, error) {
	data, _ := json.Marshal(struct {
		Size   int64 `json:"size"`
		Height int64 `json:"height"`
	}{app.state.Size, app.state.Height})
	return &types.ResponseInfo{
		Data:             string(data),
		Version:          version.ABCIVersion,
		AppVersion:       AppVersion,
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.state.Hash(),
	}, nil
}
