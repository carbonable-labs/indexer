package starknet

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/charmbracelet/log"
)

var (
	Zero = FeltFromString("0x0")
	One  = FeltFromString("0x1")
	Two  = FeltFromString("0x2")
)

func FeltFromString(s string) *felt.Felt {
	var felt felt.Felt
	err := felt.UnmarshalJSON([]byte(s))
	if err != nil {
		log.Error(err)
		return nil
	}
	return &felt
}

func FeltFromUint64(i uint64) *felt.Felt {
	var felt felt.Felt
	felt.SetUint64(i)
	return &felt
}

func HexStringToUint64(s string) uint64 {
	felt := FeltFromString(s)
	if felt == nil {
		log.Error(fmt.Sprintf("failed to parse felt from string %s", s))
		return 0
	}

	return felt.Uint64()
}
