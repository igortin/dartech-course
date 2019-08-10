package main

import (
	"dartech-course/proxy-cli/darproxy"
	"encoding/json"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

var (
	filePath string
	sep      = "/"
	count    = &darproxy.Count{}
	quit     = make(chan string)
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
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
	<-quit
}

func run(c *cli.Context) error {
	proxyConfigs := &darproxy.ProxyConfigs{}
	err := getCfg(filePath, proxyConfigs)
	if err != nil {
		return err
	}
	for _, conf := range proxyConfigs.Configs {
		cmd := darproxy.NewProxy(&http.Server{
			Addr:         conf.Port,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		}, conf)
		go cmd.Run(quit)
	}
	return nil
}

func getCfg(path string, serviceCfg *darproxy.ProxyConfigs) error {
	var b []byte
	if path == "" {
		path = os.Getenv("HOME") + sep + ".darproxy/config.json"
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	_ = json.Unmarshal(b, &serviceCfg)
	if err != nil {
		return err
	}
	return nil
}
