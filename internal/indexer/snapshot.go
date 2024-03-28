package indexer

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"slices"

	"github.com/carbonable-labs/indexer/internal/starknet"
	"github.com/carbonable-labs/indexer/internal/storage"
	"github.com/charmbracelet/log"
)

type ContractIndex struct {
	Blocks      []uint64
	LatestBlock uint64
}

func NewContractIndex(startBlock uint64) *ContractIndex {
	return &ContractIndex{
		LatestBlock: startBlock,
		Blocks:      []uint64{},
	}
}

func (c *ContractIndex) SetLatestBlock(block uint64) {
	c.LatestBlock = block
}

func (c *ContractIndex) AddBlock(block uint64) {
	if !slices.Contains(c.Blocks, block) {
		c.Blocks = append(c.Blocks, block)
	}
}

func (c *ContractIndex) Encode() (bytes.Buffer, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(c)
	if err != nil {
		return buf, err
	}

	return buf, nil
}

func (c *ContractIndex) Decode(buf []byte) error {
	decoder := gob.NewDecoder(bytes.NewBuffer(buf))
	err := decoder.Decode(c)
	if err != nil {
		return err
	}

	return nil
}

func getContractSnapshot(storage storage.Storage, address string, block uint64) (*ContractIndex, []byte) {
	contractIdxKey := []byte(fmt.Sprintf("IDX#%s", address))

	contractIdx := NewContractIndex(block)
	if storage.Has(contractIdxKey) {
		idx := storage.Get(contractIdxKey)
		err := contractIdx.Decode(idx)
		if err != nil {
			log.Error("failed to decode contract index", "error", err, "contract", address)
		}
	}
	return contractIdx, contractIdxKey
}

func replayBlocks(ctx context.Context, storage storage.Storage, blocks []uint64, blockCh chan<- starknet.GetBlockResponse) {
	for _, block := range blocks {
		if ctx.Err() != nil {
			return
		}
		block, err := fetchBlock(storage, block)
		if err != nil {
			log.Error("failed to get block", "error", err)
			continue
		}
		blockCh <- *block
	}
}

func saveContractInterestingBlock(storage storage.Storage, address string, block uint64) {
	idx, key := getContractSnapshot(storage, address, block)

	idx.AddBlock(block)

	buf, err := idx.Encode()
	if err != nil {
		log.Error("failed to encode contract index", "error", err, "contract", address)
	}
	if err := storage.Set(key, buf.Bytes()); err != nil {
		log.Error("failed to store contract index", "error", err, "contract", address)
	}
}
