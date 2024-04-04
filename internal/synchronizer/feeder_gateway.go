package synchronizer

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"

	"github.com/carbonable-labs/indexer/internal/starknet"
	"github.com/carbonable-labs/indexer/internal/storage"
	"github.com/charmbracelet/log"
)

type FeederGateway struct {
	storage storage.Storage
	client  *starknet.FeederGatewayClient
	msgch   chan *starknet.GetBlockResponse
}

func NewFeederGatewaySynchronizer(client *starknet.FeederGatewayClient, storage storage.Storage) *FeederGateway {
	return &FeederGateway{
		client:  client,
		storage: storage,
		msgch:   make(chan *starknet.GetBlockResponse),
	}
}

func (g *FeederGateway) Start(ctx context.Context) {
	log.Info("Start", "datasource", "feeder_gateway")
	s := NewSynchronizer(g.client, g.storage)
	s.Start(ctx)
}

func (g *FeederGateway) SyncBlock(ctx context.Context, block uint64) error {
	log.Info("Sync", "block", block)
	return nil
}

type Synchronizer struct {
	storage storage.Storage
	client  *starknet.FeederGatewayClient
	msgch   chan *starknet.GetBlockResponse
}

func (s *Synchronizer) Start(ctx context.Context) {
	go s.start(ctx)
	for {
		msg := <-s.msgch
		s.SyncBlock(ctx, msg)
	}
}

func (s *Synchronizer) start(ctx context.Context) {
	block := uint64(0)
	lastBlock, err := s.getLatestBlock()
	if err != nil {
		log.Error(err)
	}
	if lastBlock > block {
		block = lastBlock
	}

	for {
		resp, err := s.FetchBlock(block)
		if err != nil {
			// error will often be some timeout
			log.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}

		s.msgch <- resp
		block++
	}
}

func (s *Synchronizer) SyncBlock(ctx context.Context, block *starknet.GetBlockResponse) {
	log.Info("Sync", "block", block.BlockNumber)
	go s.storeBlock(block)

	// store data by configuration
	go s.storeLatestBlock(block.BlockNumber)
}

func (s *Synchronizer) storeBlock(block *starknet.GetBlockResponse) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(block); err != nil {
		log.Error(fmt.Sprintf("failed to encode block %s", err))
	}

	if err := s.storage.Set([]byte(fmt.Sprintf("BLOCK#%d", block.BlockNumber)), buf.Bytes()); err != nil {
		log.Error(err)
	}
}

func (s *Synchronizer) FetchBlock(blockNumber uint64) (*starknet.GetBlockResponse, error) {
	key := []byte(fmt.Sprintf("BLOCK#%d", blockNumber))
	if s.storage.Has(key) {
		block := s.storage.Get(key)
		buf := bytes.NewBuffer(block)
		decoder := gob.NewDecoder(buf)
		var resp starknet.GetBlockResponse
		err := decoder.Decode(&resp)
		if err != nil {
			log.Error(fmt.Sprintf("failed to decode block %s", err))
			return &resp, err
		}
		return &resp, nil
	}

	return s.client.GetBlock(blockNumber)
}

func (s *Synchronizer) storeLatestBlock(blockNumber uint64) {
	lastBlock, err := s.getLatestBlock()
	if err != nil {
		log.Error(err)
	}
	if lastBlock >= blockNumber {
		return
	}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(fmt.Sprintf("%d", blockNumber)); err != nil {
		log.Error(fmt.Sprintf("failed to encode block %s", err))
	}
	if err := s.storage.Set([]byte("LATEST_BLOCK"), buf.Bytes()); err != nil {
		log.Error(err)
	}
}

func (s *Synchronizer) getLatestBlock() (uint64, error) {
	res := s.storage.Get([]byte("LATEST_BLOCK"))

	buf := bytes.NewBuffer(res)
	decoder := gob.NewDecoder(buf)
	var bn string
	err := decoder.Decode(&bn)
	if err != nil {
		return 0, fmt.Errorf("failed to decode block %s", err)
	}

	num, err := strconv.ParseUint(bn, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block %s", err)
	}

	return num, nil
}

func NewSynchronizer(client *starknet.FeederGatewayClient, storage storage.Storage) *Synchronizer {
	return &Synchronizer{
		client:  client,
		storage: storage,
		msgch:   make(chan *starknet.GetBlockResponse),
	}
}
