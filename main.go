package main

import (
	"fmt"
	"log"
	"os"

	"github.com/christophercampbell/dseq/app"
	"github.com/christophercampbell/dseq/cmd"
	"github.com/urfave/cli/v2"
)

const (
	AppName = "dseq"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = AppName
	cliApp.Version = fmt.Sprintf("%v", app.AppVersion)

	cliApp.Commands = []*cli.Command{
		{
			Name:   "start",
			Usage:  "Start the node",
			Action: cmd.StartNode,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "home",
					Usage:    "Home directory `DIR`",
					Required: true,
				},
				&cli.UintFlag{
					Name:     "port",
					Usage:    "Data stream server port (6900)",
					Required: false,
					Value:    6900,
				},
			},
		}, {
			Name:   "load",
			Usage:  "Send test transaction requests",
			Action: cmd.RunLoad,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "nodes",
					Aliases:  []string{"n"},
					Usage:    "Host(&port) of nodes to send load requests, inCSV format",
					Required: true,
				},
				&cli.UintFlag{
					Name:    "requests",
					Aliases: []string{"r"},
					Usage:   "Total number of requests to send",
					Value:   10,
				},
				&cli.UintFlag{
					Name:    "concurrency",
					Aliases: []string{"c"},
					Usage:   "Number of concurrent requests",
					Value:   1,
				},
			},
		}, {
			Name:   "read",
			Usage:  "Read a data stream",
			Action: cmd.ReadStream,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "node",
					Usage:    "Node to read data stream from (host:port)",
					Required: true,
				},
				&cli.UintFlag{
					Name:     "from",
					Usage:    "Where to start the data stream from",
					Required: false,
					Value:    0,
				},
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
