package main

import (
	"encoding/json"
	"os"
)

var (
	Cmd = &darproxy{}
	// count    int
)


type darproxy struct {
	config Cfg
}

func (cmd *darproxy) GetCfg() error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&cmd.config)
	if err != nil {
		return err
	}
	return nil
}

func (cmd *darproxy) Start() error {
	HttpServer()
	return nil
}