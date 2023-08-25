package main

import (
	"fmt"
	"log"
	"os"

	"dseq/app"
	"dseq/cmd"
	"github.com/urfave/cli/v2"
)

const (
	AppName = "dseq"
)

var (
	homeFlag = cli.StringFlag{
		Name:     "home",
		Aliases:  []string{"h"},
		Usage:    "Home directory of the node `DIR`",
		Required: true,
	}
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
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
