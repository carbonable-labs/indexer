package cli_utils

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
-d, --datasource [full_node | feeder_gateway]
	Choose a datasource to use either feeder gateway (feeder_gateway) or full node (full_node) as datasource.
-h, --help prints help information
`

func CreateDatasourceFlag() DatasourceFlag {
	datasource := DatasourceFlag{}
	datasource.DatasourceBinding()

	flag.Usage = func() { log.Print(usage) }
	flag.Parse()

	return datasource
}
