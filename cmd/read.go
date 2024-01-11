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
		panic(err)
	}

	stream.FromEntry = from
	stream.SetProcessEntryFunc(printEntryNum)

	err = stream.Start()
	if err != nil {
		panic(err)
	}

	err = stream.ExecCommand(datastreamer.CmdStart)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return stream.ExecCommand(datastreamer.CmdStop)
}

func printEntryNum(e *datastreamer.FileEntry, c *datastreamer.StreamClient, s *datastreamer.StreamServer) error {
	fmt.Printf("PROCESS entry(%s): %d | %d | %d | %s\n", c.Id, e.Number, e.Length, e.Type, hexutil.Encode(e.Data))
	return nil
}
