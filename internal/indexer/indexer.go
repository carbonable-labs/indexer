package indexer

import (
	"context"
	"sync"
	"time"

	"github.com/carbonable-labs/indexer/internal/dispatcher"
	"github.com/carbonable-labs/indexer/internal/starknet"
	"github.com/carbonable-labs/indexer/internal/storage"
	"github.com/charmbracelet/log"
)

var configHashes = map[string]string{
	"app1": "hash1",
	"app2": "hash2",
}

func Run(ctx context.Context, client *starknet.FeederGatewayClient, storage storage.Storage, bus dispatcher.EventDispatcher, errCh chan<- error) {
	// get registered apps with their configuration
	// for each app create a new indexer in a goroutine
	// each goroutine will have the configuration hash as ID
	// when configuration changes, hash changes then we stop the indexer and start a new one

	var wg sync.WaitGroup
	for app, hash := range configHashes {
		wg.Add(1)
		go func(app string, hash string) {
			log.Info("Indexer started", "app", app, "hash", hash)
			i := NewIndexer(app, hash, client, storage, bus)
			err := i.Start(ctx)
			errCh <- err
			wg.Done()
		}(app, hash)
	}
	wg.Wait()
}

type Indexer struct {
	storage storage.Storage
	bus     dispatcher.EventDispatcher
	client  *starknet.FeederGatewayClient
	app     string
	hash    string
}

func NewIndexer(app string, hash string, client *starknet.FeederGatewayClient, storage storage.Storage, bus dispatcher.EventDispatcher) *Indexer {
	return &Indexer{
		app:     app,
		hash:    hash,
		client:  client,
		storage: storage,
		bus:     bus,
	}
}

func (i *Indexer) Start(ctx context.Context) error {
	for {
		time.Sleep(5 * time.Second)
		log.Debug("Indexer running", "app", i.app, "hash", i.hash)
	}
	return nil
}
