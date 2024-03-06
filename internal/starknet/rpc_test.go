package starknet_test

import (
	"testing"

	"github.com/carbonable-labs/indexer/internal/starknet"
	"gotest.tools/assert"
)

func TestCallSlotUri(t *testing.T) {
	goerli := starknet.GoerliJsonRpcStarknetClient()

	slotUri, err := starknet.GetSlotUri(goerli, "0x04b9f63c40668305ff651677f97424921bcd1b781aafa66d1b4948a87f056d0d", uint64(1))
	if err != nil {
		t.Errorf("error while testing GetSlotUri : %s", err)
	}

	assert.Equal(t, slotUri.Name, "Banegas Farm")
}

func TestCallSlotOf(t *testing.T) {
	goerli := starknet.GoerliJsonRpcStarknetClient()

	slot, err := starknet.GetSlotOf(goerli, "0x04b9f63c40668305ff651677f97424921bcd1b781aafa66d1b4948a87f056d0d", uint64(1))
	if err != nil {
		t.Errorf("error while testing GetSlotOf : %s", err)
	}

	// Token ID 1 is in slot 1
	assert.Equal(t, slot, uint64(1))
}
