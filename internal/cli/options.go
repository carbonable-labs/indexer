package cli

import (
	"errors"
	"flag"
)

var ErrInvalidDatasource = errors.New("invalid datasource")

type Datasource string

const (
	FeederGateway Datasource = "feeder_gateway"
	FullNode      Datasource = "full_node"
)

type IndexerOptions struct {
	Datasource string
	SpinUpRpc  bool
}

func (d *IndexerOptions) BindFlags() {
	flag.StringVar(&d.Datasource, "datasource", string(FeederGateway), "")
	flag.StringVar(&d.Datasource, "d", string(FeederGateway), "")
	flag.BoolVar(&d.SpinUpRpc, "nori", false, "")
	flag.BoolVar(&d.SpinUpRpc, "rpc", false, "")
}

func (d *IndexerOptions) validate() error {
	switch Datasource(d.Datasource) {
	case FullNode:
	case FeederGateway:
		return nil
	default:
		flag.Usage()
		return ErrInvalidDatasource
	}
	return nil
}
