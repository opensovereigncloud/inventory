package main

import (
	"os"

	"github.com/onmetal/inventory/pkg/app"
)

func main() {
	appInstance, ret := app.NewNICUpdaterApp()
	if ret != 0 {
		os.Exit(ret)
	}
	ret = appInstance.Run()
	os.Exit(ret)
}
