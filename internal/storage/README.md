# storage

Public interface in `storage.go`

- `Get(id []byte)`
- `Has(id []byte)`
- `Set(key []byte, value []byte)`
- `Scan(prefix []byte)`

Implementations for `nats`, `pebble` & `scylladb` present in files with the same name.
