package rpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/log"
	"github.com/keep-starknet-strange/nori"
)

func RunRpc(cancel context.CancelFunc) {
	libSig := make(chan interface{})
	config, err := setLoadLoadBalancerConfig(os.Args[1])
	if err != nil {
		log.Fatal("error reading config file", "err", err)
	}

	go runLoadBalancer(cancel, config, libSig)
}

func setLoadLoadBalancerConfig(path string) (*nori.Config, error) {
	config := new(nori.Config)
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	return config, nil
}

func runLoadBalancer(cancel context.CancelFunc, config *nori.Config, sig <-chan interface{}) {
	if config.Server.EnablePprof {
		log.Info("starting pprof", "addr", "0.0.0.0", "port", "6060")
		pprofSrv := StartPProf("0.0.0.0", 6060)
		log.Info("started pprof server", "addr", pprofSrv.Addr)
		defer func() {
			if err := pprofSrv.Close(); err != nil {
				log.Error("failed to stop pprof server", "err", err)
			}
		}()
	}

	_, shutdown, err := nori.Start(config)
	if err != nil {
		log.Fatal("error starting nori", "err", err)
	}

	select {
	case <-sig:
		shutdown()
		cancel()
	}
}

func StartPProf(hostname string, port int) *http.Server {
	mux := http.NewServeMux()

	// have to do below to support multiple servers, since the
	// pprof import only uses DefaultServeMux
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	addr := net.JoinHostPort(hostname, strconv.Itoa(port))
	srv := &http.Server{
		Handler: mux,
		Addr:    addr,
	}

	go srv.ListenAndServe()

	return srv
}
