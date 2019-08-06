package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

var (
	count int
)

const (
	index0 = iota
	index1
)

func HttpServer() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(map[string]bool{"ok": true}) })
	router.HandleFunc(Cmd.config.Upstreams[0].Path, Handler).Methods(Cmd.config.Upstreams[index0].Method)
	router.HandleFunc(Cmd.config.Upstreams[1].Path, Handler).Methods(Cmd.config.Upstreams[index1].Method)
	srv := &http.Server{
		Handler: router,
		Addr:    Cmd.config.Port,
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

func Handler(w http.ResponseWriter, r *http.Request) {
	var body []byte
	switch r.URL.Path {
		case Cmd.config.Upstreams[index0].Path:
			switch Cmd.config.Upstreams[index0].ProxyMethod {
				case "round-robin":
					count++
					body,_ = GetResponseRoundRobin(count, index0)
				case "anycast":
					body,_ = GetResponseANYCAST(index0)
			}
		case Cmd.config.Upstreams[index1].Path:
			switch Cmd.config.Upstreams[index1].ProxyMethod {
				case "round-robin":
					count++
					body, _ = GetResponseRoundRobin(count, index1)
				case "anycast":
					body,_ = GetResponseANYCAST(index1)
				}
		}
	w.Write(body)
}