package synchronizer

import (
	"context"
)

type BlockDatasource interface {
	Start(context.Context)
	SyncBlock(context.Context, uint64) error
}
