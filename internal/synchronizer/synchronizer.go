package synchronizer

import (
	"context"

	"github.com/carbonable-labs/indexer/internal/cli"
	"github.com/carbonable-labs/indexer/internal/starknet"
	"github.com/carbonable-labs/indexer/internal/storage"
	"github.com/charmbracelet/log"
)

func Run(ctx context.Context, opts *cli.IndexerOptions, client *starknet.FeederGatewayClient, storage storage.Storage) {
	var s BlockDatasource
	switch cli.Datasource(opts.Datasource) {
	case cli.FullNode:
		s = &FullNode{}
	case cli.FeederGateway:
		s = NewFeederGatewaySynchronizer(client, storage)
	default:
		log.Fatal("invalid datasource")
	}
	s.Start(ctx)
}
