package datasource

import "context"

type FullNode struct {
}

func (fg *FullNode) SyncBlock(ctx context.Context, block uint64) error {
	panic("Not implemented")
}
