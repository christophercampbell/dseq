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

func ReadStream(cli *cli.Context) error {
	node := cli.String("node")
	from := cli.Uint64("from")

	stream, err := datastreamer.NewClient(node, datastreamer.StreamType(1))
	if err != nil {
		return err
	}

	stream.FromEntry = from
	stream.SetProcessEntryFunc(printEntryNum)

	err = stream.Start()
	if err != nil {
		return err
	}

	err = stream.ExecCommand(datastreamer.CmdStart)
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return stream.ExecCommand(datastreamer.CmdStop)
}

func printEntryNum(e *datastreamer.FileEntry, c *datastreamer.StreamClient, _ *datastreamer.StreamServer) error {
	kind := "unknown"
	switch e.Type {
	case 1:
		kind = "block start"
	case 2:
		kind = "transaction"
	case 3:
		kind = "block end"
	}
	fmt.Printf("%6d | %11s | %s\n", e.Number, kind, hexutil.Encode(e.Data))
	return nil
}
