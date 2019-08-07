package main

import (
	"dartech-course/proxy-cli/darproxy"
	"github.com/urfave/cli"
	"log"
	"os"
	"sort"
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
					Destination: &darproxy.FilePath,
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

// Command to run cli darproxy
func run(c *cli.Context) error {
	darproxy.FilePath = darproxy.GetPath(darproxy.FilePath)
	err := darproxy.Cmd.GetCfg()
	if err != nil {
		return  err
	}
	err = darproxy.Cmd.Start()
	if err != nil {
		return  err
	}
	return nil
}