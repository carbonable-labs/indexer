package synchronizer

import (
	"context"
)

type BlockDatasource interface {
	SyncBlock(ctx context.Context, block uint64) error
}
