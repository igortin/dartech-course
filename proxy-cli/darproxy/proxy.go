package darproxy

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	index0 = iota
	index1
	)

type ioProxy interface {
	Start() error
	Reload() error
}

type proxy struct {
	cfg Config
}

func (cmd *proxy) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {json.NewEncoder(w).Encode(map[string]bool{"ok": true})})
	router.HandleFunc(cmd.cfg.Upstreams[0].Path, cmd.makeHandlers(index0)).Methods(cmd.cfg.Upstreams[index0].Method)
	router.HandleFunc(cmd.cfg.Upstreams[1].Path, cmd.makeHandlers(index1)).Methods(cmd.cfg.Upstreams[index1].Method)
	srv := &http.Server{
		Handler: router,
		Addr:    cmd.cfg.Port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (cmd *proxy) Reload() error {
	panic("implement me")
	return nil
}

func (cmd *proxy) Stop() error {
	panic("implement me")
	return nil
}

func NewProxy(conf Config) ioProxy {
	return &proxy{conf}
}

func (cmd *proxy) makeHandlers(index int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch cmd.cfg.Upstreams[index].ProxyMethod {
		case "round-robin":
			body,_ := GetResponseRoundRobin(index, cmd.cfg)
			w.Write(body)
		case "anycast":
			body,_ := GetResponseAnycast(index, cmd.cfg)
			w.Write(body)
		}
	}
}