package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/christophercampbell/dseq/app"
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
		return fmt.Errorf("failed to read config: %w", err)
	}
	if err = viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	if err = cfg.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	state, err := app.NewState(homeDir)
	if err != nil {
		return fmt.Errorf("failed to create state: %w", err)
	}
	defer func() {
		if closeErr := state.Close(); closeErr != nil {
			fmt.Printf("Error closing state: %v\n", closeErr)
		}
	}()

	pv := privval.LoadFilePV(
		cfg.PrivValidatorKeyFile(),
		cfg.PrivValidatorStateFile(),
	)

	addr := common.BytesToAddress(pv.GetAddress().Bytes())

	var nodeKey *p2p.NodeKey
	if nodeKey, err = p2p.LoadNodeKey(cfg.NodeKeyFile()); err != nil {
		return fmt.Errorf("failed to load node key: %w", err)
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
		return fmt.Errorf("failed to create stream server: %w", err)
	}

	if err = streamServer.Start(); err != nil {
		return fmt.Errorf("failed to start stream server: %w", err)
	}

	logger := cmtlog.NewTMLogger(cmtlog.NewSyncWriter(os.Stdout))
	if logger, err = flags.ParseLogLevel(cfg.LogLevel, logger, config.DefaultLogLevel); err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	sequencer, err := app.NewSequencer(
		logger,
		app.WithIdentity(cfg.Moniker),
		app.WithAddress(addr),
		app.WithState(state),
		app.WithDataServer(streamServer),
	)
	if err != nil {
		return fmt.Errorf("failed to create sequencer: %w", err)
	}

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
		return fmt.Errorf("failed to create node: %w", err)
	}

	if err = n.Start(); err != nil {
		return fmt.Errorf("failed to start node: %w", err)
	}

	defer func() {
		if stopErr := n.Stop(); stopErr != nil {
			fmt.Printf("Error stopping node: %v\n", stopErr)
		}
		n.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}
