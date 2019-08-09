package darproxy

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	index0 = iota
	index1
)

type ioProxy interface {
	Start() error
	Reload(Config) error
}

type proxy struct {
	cfg Config
}

func (cmd *proxy) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(map[string]bool{"ok": true}) })
	router.HandleFunc(cmd.cfg.Upstreams[0].Path, cmd.makeHandlers(index0)).Methods(cmd.cfg.Upstreams[index0].Method)
	router.HandleFunc(cmd.cfg.Upstreams[1].Path, cmd.makeHandlers(index1)).Methods(cmd.cfg.Upstreams[index1].Method)
	srv := &http.Server{
		Handler: router,
		Addr:    cmd.cfg.Port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}



	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	log.Println("server started")
	stop := make(chan os.Signal, 1)							// create ch for os signal
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)		// configure wait signals good channel
	<-stop													// until no object good in channel block below code lines
	// below code
	log.Println("received stop signal")				// if get from good channel obj go exec below code

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("call to shutdown")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Println("Server stopped")
	}
	time.Sleep(5 * time.Second)
	return nil
}

func (cmd *proxy) Reload(config Config) error {
	cmd.cfg = config
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
			body, _ := GetResponseRoundRobin(index, cmd.cfg)
			w.Write(body)
		case "anycast":
			body, _ := GetResponseAnycast(index, cmd.cfg)
			w.Write(body)
		}
	}
}