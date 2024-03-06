package storage

type (
	Storage interface {
		Get(id []byte) []byte
		Has(id []byte) bool
		Set(key []byte, value []byte) error
		Scan(prefix []byte) [][]byte
	}
)
