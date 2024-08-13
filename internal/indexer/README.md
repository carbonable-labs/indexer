# indexer

This package is responsible of indexing events and txs based on the configuration input.

Package entrypoint : `Run`

This function checks for registered configurations and starts an indexing process for each of those configurations.

- Configurations are based on the `app_name` field and are mapped to a hash of their content.
- When a configuration is updated - the hash for the name changes - the indexer for that configuration is updated. You don't loose previous configuration hash but it will stop indexing for new events / txs.
- An indexer subprocess starts at block `starting_block` and pull data from the `synchronizer` process.
- For each block, the indexer checks if their are events / txs related to the given configuration
- When data related to the configuration is found, events are dispatched through `dispatcher` package and block is stored as related - when you restart this indexer it will only pickup blocks where related data exists.

## `indexer.go`

All indexing logic

## `snapshot.go`

Code to store configuration snapshot (the blocks where contract have data related from the contract to index.)

## `blocks.go`

Helper functions to fetch blocks from storage and stream them into a channel.

- `iterateBlocks`
- `fetchBlock`
- `addMissingBlock`

## `config.go`

Helper functions to check for configuration changes and fetch configurations from storage.

- `checkConfigChange`
- `fetchConfigurations`
- `getConfigurationDiffs`
