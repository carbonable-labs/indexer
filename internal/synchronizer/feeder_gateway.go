package synchronizer

import (
	"context"

	"github.com/charmbracelet/log"
)

type FeederGateway struct {
}

func (fg *FeederGateway) SyncBlock(ctx context.Context, block uint64) error {
	log.Info("Sync", "block", block)
	return nil
}
