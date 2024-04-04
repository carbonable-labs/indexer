package cli

import (
	"flag"
	"log"
)

/*
Specification of various CLI parameters.
Update when a new CLI parameter is added.
*/
const usage = `
Command Line Options:
	-d, --datasource [full_node | feeder_gateway]      Choose a datasource to use either feeder gateway (feeder_gateway) or full node (full_node) as datasource.
	-nori, --rpc                                       Spin up rpc server for rpc load balancing
	-h, --help                                         Prints help information

`

func ConfigureOptions() (*IndexerOptions, error) {
	opts := &IndexerOptions{}
	opts.BindFlags()

	flag.Usage = func() { log.Print(usage) }
	flag.Parse()

	err := opts.validate()

	return opts, err
}
