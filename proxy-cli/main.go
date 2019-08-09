package main

import (
	"dartech-course/proxy-cli/darproxy"
	"encoding/json"
	"github.com/urfave/cli"
	"log"
	"os"
	"sort"
)

var (
	filePath string
	sep      = "/"
	count = &darproxy.Count{}
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
					Name:        "config, c",
					Usage:       "Load config from `FILE`",
					Destination: &filePath,
				},
			},
			Action: run,
		},
		{
			Name:      "reload",
			Usage:     "reload service",
			UsageText: "dar-proxy reload",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "config, c",
					Usage:       "Load config from `FILE`",
					Destination: &filePath,
				},
			},
			Action: reload,
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
	conf := &darproxy.Config{}
	err := getCfg(filePath, conf)
	if err != nil {
		return err
	}
	cmd := darproxy.NewProxy(*conf)
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func reload(c *cli.Context) error {
	conf := &darproxy.Config{}
	getCfg(filePath, conf)
	panic("implement me")
	return nil
}

func getCfg(path string, config *darproxy.Config) (error) {
	if path == "" {
		path = os.Getenv("HOME") + sep + ".darproxy/config.json"
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(config)
	if err != nil {
		return err
	}
	return err
}