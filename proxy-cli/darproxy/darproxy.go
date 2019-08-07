package darproxy

import (
	"encoding/json"
	"os"
)

var (
	Cmd = &darproxy{}
	FilePath string
)
// darproxy cli Struct
type darproxy struct {
	config Cfg
}

// Method darproxy find/get cfg
func (cmd *darproxy) GetCfg() error {
	f, err := os.Open(FilePath)
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

// Method darproxy start
func (cmd *darproxy) Start() error {
	HttpServer()
	return nil
}