package starknet

import (
	"fmt"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"golang.org/x/crypto/sha3"
)

// EnsureStarkFelt adds 0 padding to left side of felt
func EnsureStarkFelt(felt string) string {
	return fmt.Sprintf("0x%064s", strings.Replace(felt, "0x", "", 1))
}

// encode storage variable name to felt
func StarknetKeccak(b []byte) (*felt.Felt, error) {
	h := sha3.NewLegacyKeccak256()
	_, err := h.Write(b)
	if err != nil {
		return nil, err
	}
	d := h.Sum(nil)
	// Remove the first 6 bits from the first byte
	d[0] &= 3
	return new(felt.Felt).SetBytes(d), nil
}
