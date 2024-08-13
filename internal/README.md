# internal

Internal go packages that to run the application.

- `api` - api endpoints to register configurations
- `cli` - cli options + configuration of the application through command line
- `config` - configuration building blocks of the application. Responsible for parsing configuration files and building configuration models.
- `dispatcher` - behaviour for event dispatching. Plug anything that you want to use to dispatch messages.
- `indexer` - main package that runs indexer logic.
- `rpc` - related to run starknet rpc load balancer using nori
- `starknet` - starknet specific code
- `storage` - behaviour for storage.
- `synchronizer` - synchronizer is in charge of synchronizing the indexer with the sequencer gateway.
