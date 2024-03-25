package cli_utils

import (
	"flag"
	"fmt"

	"github.com/carbonable-labs/indexer/internal/synchronizer"
)

type Datasource string

const (
	FeederGateway Datasource = "feeder_gateway"
	FullNode      Datasource = "full_node"
)

type DatasourceFlag struct {
	Value    string
	Is_valid bool
}

func (d *DatasourceFlag) DatasourceBinding() {
	flag.StringVar(&d.Value, "datasource", string(FeederGateway), "")
	flag.StringVar(&d.Value, "d", string(FeederGateway), "")
}

func (d *DatasourceFlag) Validation() synchronizer.BlockDatasource {
	switch Datasource(d.Value) {
	case FullNode:
		return &synchronizer.FullNode{}
	case FeederGateway:
		return &synchronizer.FeederGateway{}
	default:
		flag.Usage()
		fmt.Print("Error: Value proportioned for datasource is invalid.\n")
		return nil
	}
}
