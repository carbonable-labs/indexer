package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/carbonable-labs/indexer/internal/api"
	"github.com/carbonable-labs/indexer/internal/cli"
	"github.com/carbonable-labs/indexer/internal/dispatcher"
	"github.com/carbonable-labs/indexer/internal/indexer"
	"github.com/carbonable-labs/indexer/internal/starknet"
	"github.com/carbonable-labs/indexer/internal/storage"
	"github.com/carbonable-labs/indexer/internal/synchronizer"
	"github.com/charmbracelet/log"
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
	ctx, cancel := context.WithCancel(context.Background())

	setupLogLevel()
	setupCancel(cancel)

	// Setup of --datasource, -d CLI flag
	opts, err := cli.ConfigureOptions()
	if err != nil {
		log.Debug("error configuring cli options", "err", err)
		os.Exit(1)
	}

	natsUrl := os.Getenv("NATS_URL")
	natsToken := os.Getenv("NATS_TOKEN")

	client := starknet.NewSepoliaFeederGatewayClient()
	storage := storage.NewNatsStorage(storage.WithUrl(natsUrl), storage.WithToken(natsToken))
	dispatcher := dispatcher.NewNatsDispatcher(dispatcher.WithUrl(natsUrl), dispatcher.WithToken(natsToken))

	go synchronizer.Run(ctx, opts, client, storage)
	go api.Run(ctx, storage)
	go indexer.Run(ctx, client, storage, dispatcher)

	// FIX: dependencies version clash with nori works well as a standalone
	// if opts.SpinUpRpc {
	// 	go rpc.RunRpc(cancel)
	// }

	<-ctx.Done()
	log.Info("shutting down gracefully")
}

func setupCancel(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()
}

func setupLogLevel() {
	logLevelString := os.Getenv("LOG_LEVEL")
	var logLevel log.Level
	switch logLevelString {
	case "debug":
		logLevel = log.DebugLevel
	case "info":
		logLevel = log.InfoLevel
	case "warn":
		logLevel = log.WarnLevel
	case "error":
		logLevel = log.ErrorLevel
	default:
		logLevel = log.InfoLevel
	}
	_ = logLevel

	log.SetLevel(logLevel)
}
