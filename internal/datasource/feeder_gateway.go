package datasource

import "context"

type FeederGateway struct {
}

func (fg *FeederGateway) SyncBlock(ctx context.Context, block uint64) error {
	return nil
}
