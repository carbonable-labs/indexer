# config

Application configuration. Parses json/yaml into a configuration model as so

```json
{
  "app_name": "carbonable-indexer",
  "start_block": 1,
  "contracts": [
    {
      "name": "carbonable-token",
      "address": "0x06a09ccb1caaecf3d9683efe335a667b2169a409a5293e49c23cb8b539c4e48",
      "events": {
        "Transfer": "token:transfer"
      }
    }
  ]
}
```

- `app_name` is the name of the application.
- `start_block` is the block number to start indexing from. **must be > 0**
- `contracts` is an array of contracts to index.
  - `name` is the name of the contract
  - `address` is the address of the contract
  - `events` is an object of events to listen to.
    - key is the event name as defined in the contract
    - value is the event name you want to receive from the indexer.
