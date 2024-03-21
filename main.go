package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/carbonable-labs/indexer/internal/starknet"
	"github.com/carbonable-labs/indexer/internal/storage"
	"github.com/carbonable-labs/indexer/internal/synchronizer"
	"golang.org/x/sync/errgroup"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/log"
	"github.com/keep-starknet-strange/nori"
	"golang.org/x/exp/slog"
)

const welcomeMessage = "Sheshat ... Indexing"

var (
	GitVersion = ""
	GitCommit  = ""
	GitDate    = ""
)

// starting block
// indexing configuration

// get all contracts declared at in the genesis dataget_class_by_hash
// -> each time config is changed, know where to start indexing from
// -> keep all indexing configurations to enable fast retrieval / per contract

// we may want to ignore what is before starting block as it is not required to have data
// for each contract -> index all events in a event driven way
// store txs and state updates as customer may want to retrieve data based on this.

// First step
// fetch all data related to contracts
// store data

// Second step
// stream data into message broker

// Every reload or reindex is based off a last_event_id (ulid based on time)
// replayed from database to get faster indexing

func main() {
	ctx, _ := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	// Set up logger with a default INFO level in case we fail to parse flags.
	// Otherwise the final critical log won't show what the parsing error was.
	log.SetDefault(log.NewLogger(log.LogfmtHandlerWithLevel(os.Stdout, log.LvlInfo)))

	log.Info("starting nori", "version", GitVersion, "commit", GitCommit, "date", GitDate)

	if len(os.Args) < 2 {
		log.Crit("must specify a config file on the command line")
	}

	config, err := loadLoadBalancerConfig(os.Args[1])
	if err != nil {
		log.Crit("error reading config file", "err", err)
	}

	// update log level from config
	logLevelString := config.Server.LogLevel
	var logLevel slog.Level
	switch logLevelString {
	case "trace":
		logLevel = log.LevelTrace
	case "debug":
		logLevel = log.LevelDebug
	case "info":
		logLevel = log.LevelInfo
	case "warn":
		logLevel = log.LevelWarn
	case "error":
		logLevel = log.LevelError
	case "crit":
		logLevel = log.LevelCrit
	default:
		logLevel = log.LevelInfo
		log.Warn("invalid server.log_level set: " + logLevelString)
	}
	log.SetDefault(log.NewLogger(log.LogfmtHandlerWithLevel(os.Stdout, logLevel)))

	g.Go(func() error {
		return runSynchronizer(ctx)
	})

	g.Go(func() error {
		return runLoadBalancer(ctx, config)
	})

	if err := g.Wait(); err != nil {
		// defer cancel()
		log.Crit(fmt.Sprintf("Indexer process terminated: %s", err))
	}
}

func loadLoadBalancerConfig(path string) (*nori.Config, error) {
	config := new(nori.Config)
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	return config, nil
}

func runSynchronizer(ctx context.Context) error {
	fmt.Println(welcomeMessage)

	errCh := make(chan error)
	client := starknet.NewSepoliaFeederGatewayClient()
	storage := storage.NewPebbleStorage()

	go synchronizer.Run(ctx, client, storage, errCh)

	err := <-errCh
	return fmt.Errorf("error while syncing network: %s", err)
}

func runLoadBalancer(_ context.Context, config *nori.Config) error {
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
		log.Crit("error starting nori", "err", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	recvSig := <-sig
	log.Info("caught signal, shutting down", "signal", recvSig)
	shutdown()

	return nil
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
