package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/urfave/cli/v2"
)

// ReadStream reads and displays data from a stream starting from the specified entry.
func ReadStream(cli *cli.Context) error {
	node := cli.String("node")
	from := cli.Uint64("from")

	stream, err := datastreamer.NewClient(node, datastreamer.StreamType(1))
	if err != nil {
		return fmt.Errorf("failed to create stream client: %w", err)
	}

	stream.FromEntry = from
	stream.SetProcessEntryFunc(printEntryNum)

	if err := stream.Start(); err != nil {
		return fmt.Errorf("failed to start stream: %w", err)
	}

	if err := stream.ExecCommand(datastreamer.CmdStart); err != nil {
		return fmt.Errorf("failed to execute start command: %w", err)
	}

	// Set up signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	if err := stream.ExecCommand(datastreamer.CmdStop); err != nil {
		return fmt.Errorf("failed to execute stop command: %w", err)
	}

	return nil
}

// printEntryNum prints information about a stream entry.
func printEntryNum(e *datastreamer.FileEntry, c *datastreamer.StreamClient, _ *datastreamer.StreamServer) error {
	if e == nil {
		return fmt.Errorf("received nil entry")
	}

	kind := "unknown"
	switch e.Type {
	case 1:
		kind = "block start"
	case 2:
		kind = "transaction"
	case 3:
		kind = "block end"
	default:
		kind = fmt.Sprintf("unknown type %d", e.Type)
	}

	fmt.Printf("%6d | %11s | %s\n", e.Number, kind, hexutil.Encode(e.Data))
	return nil
}
