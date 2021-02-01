package main

import (
	"os"

	"github.com/onmetal/inventory/pkg/inventory"
)

func main() {
	is, ret := inventory.NewSvc()
	if ret != 0 {
		os.Exit(ret)
	}
	ret = is.Inventorize()
	os.Exit(ret)
}
