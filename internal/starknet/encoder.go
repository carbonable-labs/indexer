package starknet

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/charmbracelet/log"
)

func DecodeGob[T any](data []byte) (*T, error) {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	var resp T
	err := decoder.Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func EncodeGob[T any](data *T) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeSlice[T any](data [][]byte) ([]*T, error) {
	ds := make([]*T, len(data))
	for i, d := range data {
		r, err := DecodeGob[T](d)
		if err != nil {
			return nil, err
		}
		ds[i] = r
	}
	return ds, nil
}

func DecodeResponseToStruct[T any](data []felt.Felt, s T) error {
	// INFO: Trick here is to remove the first element as it stands for array length
	// and second as it stands for encoding "data:application/json,"
	strVal := FeltArrToBytesArr(data[2:])

	err := json.Unmarshal(strVal, &s)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
