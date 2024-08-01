package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
)

type (
	Storage interface {
		Get(id []byte) []byte
		Has(id []byte) bool
		Set(key []byte, value []byte) error
		Scan(prefix []byte) [][]byte
	}
)

func GetLatestBlock(s Storage) (uint64, error) {
	res := s.Get([]byte("latest_block"))

	buf := bytes.NewBuffer(res)
	decoder := gob.NewDecoder(buf)
	var bn string
	err := decoder.Decode(&bn)
	if err != nil {
		return 0, fmt.Errorf("failed to decode block %s", err)
	}

	num, err := strconv.ParseUint(bn, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block %s", err)
	}

	return num, nil
}
