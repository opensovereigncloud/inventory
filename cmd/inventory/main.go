package main

import (
	"os"

	"github.com/onmetal/inventory/pkg/gatherer"
)

func main() {
	is, ret := gatherer.NewSvc()
	if ret != 0 {
		os.Exit(ret)
	}
	ret = is.Gather()
	os.Exit(ret)
}
