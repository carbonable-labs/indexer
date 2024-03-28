package indexer

import (
	"context"
	"time"

	"github.com/carbonable-labs/indexer/internal/config"
	"github.com/carbonable-labs/indexer/internal/storage"
	"github.com/charmbracelet/log"
)

func checkConfigChange(ctx context.Context, storage storage.Storage, cfg config.Config, cfgCh chan<- config.Config) {
	cr := config.NewPebbleContractRepository(storage)
	for {
		if ctx.Err() != nil {
			return
		}
		time.Sleep(5 * time.Second)
		c, err := cr.GetConfig(cfg.AppName)
		if err != nil {
			log.Error("failed to get config", "error", err)
			continue
		}
		if c.Hash != cfg.Hash {
			cfgCh <- *c
		}
	}
}
