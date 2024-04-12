package indexer

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/carbonable-labs/indexer/internal/starknet"
	"github.com/carbonable-labs/indexer/internal/storage"
	"github.com/charmbracelet/log"
)

// Iterate over blocks that are sync with datasource
// Streams blocks into channel
func iterateBlocks(ctx context.Context, storage storage.Storage, block uint64, blockCh chan starknet.GetBlockResponse) {
	for {
		if ctx.Err() != nil {
			return
		}
		resp, err := fetchBlock(storage, block)
		if err != nil {
			log.Debug("failed to get block", "error", err, "block", block)

			time.Sleep(10 * time.Second)
			continue
		}

		select {
		case <-ctx.Done():
			return
		case blockCh <- *resp:
			block++
		}
	}
}

// Single piece of code to retrieve block from storage with specific key
func fetchBlock(storage storage.Storage, block uint64) (*starknet.GetBlockResponse, error) {
	key := []byte(fmt.Sprintf("block.%d", block))
	if storage.Has(key) {
		block := storage.Get(key)
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
	return nil, errors.New("block not found")
}
