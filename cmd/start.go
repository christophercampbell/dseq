package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"dseq/app"
	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/libs/cli/flags"
	cmtlog "github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func StartNode(cli *cli.Context) error {

	homeDir := cli.String("home")
	dataPort := uint16(cli.Uint("port"))

	cfg := config.DefaultConfig()
	cfg.SetRoot(homeDir)

	var err error
	viper.SetConfigFile(fmt.Sprintf("%s/%s", homeDir, "config/config.toml"))
	if err = viper.ReadInConfig(); err != nil {
		return err
	}
	if err = viper.Unmarshal(cfg); err != nil {
		return err
	}
	if err = cfg.ValidateBasic(); err != nil {
		return err
	}

	state := app.NewState(homeDir)
	defer state.Close()

	pv := privval.LoadFilePV(
		cfg.PrivValidatorKeyFile(),
		cfg.PrivValidatorStateFile(),
	)

	addr := common.BytesToAddress(pv.GetAddress().Bytes())

	var nodeKey *p2p.NodeKey
	if nodeKey, err = p2p.LoadNodeKey(cfg.NodeKeyFile()); err != nil {
		return err
	}

	streamFile := strings.Join([]string{homeDir, "dseq.bin"}, "/")

	streamServer, err := datastreamer.NewServer(
		dataPort,
		1,
		1,
		datastreamer.StreamType(app.StSequencer),
		streamFile,
		nil,
	)
	if err != nil {
		panic(err)
	}

	err = streamServer.Start()
	if err != nil {
		panic(err)
	}

	logger := cmtlog.NewTMLogger(cmtlog.NewSyncWriter(os.Stdout))
	if logger, err = flags.ParseLogLevel(cfg.LogLevel, logger, config.DefaultLogLevel); err != nil {
		return err
	}

	sequencer := app.NewSequencer(logger, cfg.Moniker, addr, state, streamServer)

	var n *node.Node
	if n, err = node.NewNode(
		cfg,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(sequencer),
		node.DefaultGenesisDocProviderFunc(cfg),
		config.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		logger); err != nil {
		return err
	}

	if err = n.Start(); err != nil {
		return err
	}

	defer func() {
		_ = n.Stop()
		n.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return err
}
