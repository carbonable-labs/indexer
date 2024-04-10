package starknet_test

import (
	"testing"

	"github.com/carbonable-labs/indexer/internal/starknet"
	"github.com/test-go/testify/assert"
)

func TestFeltFromUint(t *testing.T) {
	testCases := []struct {
		name     string
		expected uint64
	}{
		{
			name:     "3",
			expected: uint64(3),
		},
		{
			name:     "15",
			expected: uint64(15),
		},
		{
			name:     "118",
			expected: uint64(118),
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			felt := starknet.FeltFromUint64(tc.expected)
			if nil == felt {
				t.Error("failed to create felt from uint64")
			}
			assert.Equal(t, felt.Uint64(), tc.expected)
		})
	}
}

func TestFeltFromString(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected uint64
	}{
		{
			name:     "3",
			value:    "0x3",
			expected: uint64(3),
		},
		{
			name:     "15",
			value:    "0xf",
			expected: uint64(15),
		},
		{
			name:     "118",
			value:    "0x76",
			expected: uint64(118),
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			felt := starknet.FeltFromString(tc.value)
			if nil == felt {
				t.Error("failed to create felt from uint64")
			}
			assert.Equal(t, felt.Uint64(), tc.expected)
		})
	}
}
