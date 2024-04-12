package indexer

import (
	"testing"

	"github.com/carbonable-labs/indexer/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigDiff(t *testing.T) {
	cfgs := []config.Config{}

	cfgs2 := []config.Config{
		{AppName: "app1", Hash: "hash1"},
	}

	diff := getConfigurationDiffs(cfgs, cfgs2)
	assert.Equal(t, len(diff), 1)
	assert.Equal(t, diff[0].AppName, "app1")
}

func TestConfigDiffSame(t *testing.T) {
	cfgs := []config.Config{
		{AppName: "app1", Hash: "hash1"},
	}

	cfgs2 := []config.Config{
		{AppName: "app1", Hash: "hash1"},
	}

	diff := getConfigurationDiffs(cfgs, cfgs2)
	assert.Equal(t, len(diff), 0)
}

func TestConfigDiffSameWithDiff(t *testing.T) {
	cfgs := []config.Config{
		{AppName: "app1", Hash: "hash1"},
	}

	cfgs2 := []config.Config{
		{AppName: "app1", Hash: "hash1"},
		{AppName: "app2", Hash: "hash2"},
	}

	diff := getConfigurationDiffs(cfgs, cfgs2)
	assert.Equal(t, len(diff), 1)
	assert.Equal(t, diff[0].AppName, "app2")
}
