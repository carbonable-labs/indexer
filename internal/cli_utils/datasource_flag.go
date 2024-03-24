package cli_utils

import (
	"flag"
	"fmt"
	"slices"
)

type datasource_flag struct {
	Value    string
	Is_valid bool
}

func (d *datasource_flag) datasource_binding() {
	flag.StringVar(&d.Value, "datasource", "fg", "")
	flag.StringVar(&d.Value, "d", "fg", "")
}

func (d *datasource_flag) Validation() {
	/*
		fg = feed gateway
		fn = full node rpc
	*/
	valid_values := []string{"fg", "fn"}
	d.Is_valid = slices.Contains(valid_values, d.Value)
	if !d.Is_valid {
		flag.Usage()
		fmt.Print("Error: Value proportioned for datasource is invalid.\n")
		//something to abort execution
	}
}
