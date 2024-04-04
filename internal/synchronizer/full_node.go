package synchronizer

import "context"

type FullNode struct{}

func (n *FullNode) Start(ctx context.Context) {
	panic("Not implemented")
}

func (n *FullNode) SyncBlock(ctx context.Context, block uint64) error {
	panic("Not implemented")
}
