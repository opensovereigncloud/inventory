package main

import (
	"os"

	"github.com/onmetal/inventory/pkg/gatherer"
)

func main() {
	app, ret := gatherer.NewNICUpdaterSvc()
	if ret != 0 {
		os.Exit(ret)
	}
	ret = app.Run()
	os.Exit(ret)
}
