package indexer

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"sync"
	"time"

	"github.com/carbonable-labs/indexer/internal/config"
	"github.com/carbonable-labs/indexer/internal/dispatcher"
	"github.com/carbonable-labs/indexer/internal/starknet"
	"github.com/carbonable-labs/indexer/internal/storage"
	"github.com/charmbracelet/log"
)

// Run an indexer for each app registered
func Run(ctx context.Context, client *starknet.FeederGatewayClient, storage storage.Storage, bus dispatcher.EventDispatcher, errCh chan<- error) {
	// get registered apps with their configuration
	// for each app create a new indexer in a goroutine
	// each goroutine will have the configuration hash as ID
	// when configuration changes, hash changes then we stop the indexer and start a new one

	cr := config.NewPebbleContractRepository(storage)
	configs, err := cr.GetConfigs()
	if err != nil {
		log.Fatal("unable to retrieve configurations", "error", err)
	}

	var wg sync.WaitGroup
	for _, c := range configs {
		wg.Add(1)
		go func(c config.Config) {
			log.Info("Indexer started", "app", c.AppName, "hash", c.Hash)
			i := NewIndexer(c, client, storage, bus)
			err := i.Start(ctx)
			errCh <- err
			wg.Done()
		}(c)
	}
	wg.Wait()
}

// Single contract indexer
type Indexer struct {
	storage storage.Storage
	bus     dispatcher.EventDispatcher
	client  *starknet.FeederGatewayClient
	config  config.Config
}

func NewIndexer(config config.Config, client *starknet.FeederGatewayClient, storage storage.Storage, bus dispatcher.EventDispatcher) *Indexer {
	return &Indexer{
		config:  config,
		client:  client,
		storage: storage,
		bus:     bus,
	}
}

func (i *Indexer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	blockCh := make(chan starknet.GetBlockResponse)
	cfgCh := make(chan config.Config)
	go checkConfigChange(ctx, i.storage, i.config, cfgCh)

	indexProcess(ctx, i, i.config, blockCh)

	for {
		select {
		case block := <-blockCh:
			for _, c := range i.config.Contracts {
				err := i.Index(c, block)
				if err != nil {
					log.Error("failed to index block", "error", err)
				}
			}
		case cfg := <-cfgCh:
			i.config = cfg
			log.Info("config changed", "app", cfg.AppName, "hash", cfg.Hash)
			cancel()

			return i.Start(context.Background())
		}
	}
}

func (i *Indexer) Index(contract config.Contract, block starknet.GetBlockResponse) error {
	// store txs
	// store events
	// save contract snapshot if block has data about specific config
	return nil
}

func (i *Indexer) indexTransaction(address string, block *starknet.GetBlockResponse) {
	for _, tx := range block.Transactions {
		if starknet.EnsureStarkFelt(tx.SenderAddress) != starknet.EnsureStarkFelt(address) {
			continue
		}
		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		err := encoder.Encode(tx)
		if err != nil {
			log.Error("failed to encode event", "error", err)
		}

		if err := i.storage.Set([]byte(fmt.Sprintf("%s.tx.%s", address, tx.TransactionHash)), buf.Bytes()); err != nil {
			log.Error("failed to store event", "error", err)
		}
		log.Info("Indexing tx for address", "address", address, "tx", tx.TransactionHash)

		saveContractInterestingBlock(i.storage, address, block.BlockNumber)
	}
}

func (i *Indexer) indexEvent(address string, block *starknet.GetBlockResponse) {
	for _, tx := range block.TransactionReceipts {
		for eventIdx, event := range tx.Events {
			if starknet.EnsureStarkFelt(event.FromAddress) != starknet.EnsureStarkFelt(address) {
				continue
			}

			var buf bytes.Buffer
			encoder := gob.NewEncoder(&buf)

			// Aggregating event_id to event
			eventId := fmt.Sprintf("%s_%d", tx.TransactionHash, eventIdx)
			event.EventId = eventId
			event.RecordedAt = time.Unix(int64(block.Timestamp), 0)

			err := encoder.Encode(event)
			if err != nil {
				log.Error("failed to encode event", "error", err)
			}

			if err := i.storage.Set([]byte(fmt.Sprintf("%s.event.%s", address, eventId)), buf.Bytes()); err != nil {
				log.Error("failed to store event", "error", err)
			}
			if err := i.storage.Set([]byte(fmt.Sprintf("event.%s", eventId)), buf.Bytes()); err != nil {
				log.Error("failed to store event", "error", err)
			}
			// i.nats.Publish("event:published", []byte(eventId))
			log.Info("Indexing event for address", "address", address, "eventId", eventId)

			saveContractInterestingBlock(i.storage, address, block.BlockNumber)
		}
	}
}

// Start processes to check config and fetch blocks per config
func indexProcess(ctx context.Context, i *Indexer, c config.Config, blockCh chan starknet.GetBlockResponse) {
	// check starting block
	startBlock := c.StartBlock
	for _, c := range c.Contracts {
		idx, _ := getContractSnapshot(i.storage, c.Address, startBlock)

		go replayBlocks(ctx, i.storage, idx.Blocks, blockCh)

		block := startBlock
		if idx.LatestBlock > block {
			block = idx.LatestBlock
		}

		go iterateBlocks(ctx, i.storage, block, blockCh)
	}
}
