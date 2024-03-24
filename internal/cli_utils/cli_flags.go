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
-d, --datasource [fg | fn]
	Choose a datasource to use either feeder gateway (fd) or full node (fn) as datasource.
-h, --help prints help information
`

func DatasourceFlag() datasource_flag {
	datasource := datasource_flag{}
	datasource.datasource_binding()

	flag.Usage = func() { log.Print(usage) }
	flag.Parse()

	return datasource
}
