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

func fetchConfigurations(cr config.ContractRepository, cfgChan chan<- []config.Config) {
	for {
		time.Sleep(5 * time.Second)
		cfgs, err := cr.GetConfigs()
		if err != nil {
			log.Error("failed to get configs", "error", err)
			continue
		}
		cfgChan <- cfgs
	}
}

func getConfigurationDiffs(old, new []config.Config) []config.Config {
	var diff []config.Config
	for _, n := range new {
		found := false
		for _, o := range old {
			// NOTE: we check on appName variations since app configuration reloads are managed within single indexers
			if n.AppName == o.AppName {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, n)
		}
	}
	return diff
}
