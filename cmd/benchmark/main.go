package main

import (
	"os"

	"github.com/onmetal/inventory/pkg/benchmark"
)

func main() {
	is, ret := benchmark.NewSvc()
	if ret != 0 {
		os.Exit(ret)
	}
	ret = is.Gather()
	os.Exit(ret)
}
