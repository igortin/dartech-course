package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"sort"
)

var (
	sep = "/"
	filepath string
)

func main() {
	app := cli.NewApp()
	app.Name = "darproxy"
	app.Version = "0.0.0"
	app.Usage = "darproxy run --config 999"
	app.Description = "Proxy service with different politics"
	app.Commands = []cli.Command{
		{
			Name:      "run",
			Usage:     "start proxy service",
			UsageText: "dar-proxy run",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Load config from `FILE`",
					Destination: &filepath,
				},
			},
			Action: run,
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}


func run(c *cli.Context) error {
	filepath = GetPath(filepath)
	err := Cmd.GetCfg()
	if err != nil {
		return  err
	}
	err = Cmd.Start()
	if err != nil {
		return  err
	}
	return nil
}