package darproxy

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"
)

const (
	index0 = iota
	index1
)

type ioProxy interface {
	Run() error
	Shutdown() error
}

type proxy struct {
	service         *http.Server
	stopped         bool
	router          *mux.Router
	gracefulTimeout time.Duration
	cfg             Config
}

func (cmd *proxy) Run() error {
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		if err := cmd.Shutdown(); err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			log.Println("Server stopped")
		}
	}()
	cmd.stopped = false
	cmd.router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {json.NewEncoder(w).Encode(map[string]bool{"ok": true})})
	cmd.router.HandleFunc(cmd.cfg.Upstreams[0].Path, cmd.makeHandlers(index0)).Methods(cmd.cfg.Upstreams[index0].Method)
	cmd.router.HandleFunc(cmd.cfg.Upstreams[1].Path, cmd.makeHandlers(index1)).Methods(cmd.cfg.Upstreams[index1].Method)
	return cmd.service.ListenAndServe()
}

func NewProxy(srv *http.Server, config Config) ioProxy {
	router := mux.NewRouter()
	srv.Handler = router
	graceTimeout := 5 * time.Second
	return &proxy{
		srv,
		true,
		router,
		graceTimeout,
		config,
	}
}

func (cmd *proxy) makeHandlers(index int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cmd.stopped {
			w.WriteHeader(503)
			return
		}
		switch cmd.cfg.Upstreams[index].ProxyMethod {
		case "round-robin":
			body, _ := GetResponseRoundRobin(index, cmd.cfg)
			w.Write(body)

		case "anycast":
			body, _ := GetResponseAnycast(index, cmd.cfg)
			w.Write(body)
		}
	}
}

func (cmd *proxy) Shutdown() error {
	cmd.stopped = true
	ctx, cancel := context.WithTimeout(context.Background(), cmd.gracefulTimeout)
	defer cancel()
	time.Sleep(cmd.gracefulTimeout)
	return cmd.service.Shutdown(ctx)
}