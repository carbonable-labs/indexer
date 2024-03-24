package cli_utils

import (
	"flag"
	"fmt"
	"slices"
)

type datasource_flag struct {
	Value    string
	is_valid bool
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
	d.is_valid = slices.Contains(valid_values, d.Value)
	if !d.is_valid {
		flag.Usage()
		fmt.Print("Error: Value proportioned for datasource is invalid.\n")
		//something to abort execution
	}
}
