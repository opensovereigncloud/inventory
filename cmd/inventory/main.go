package main

import (
	"github.com/onmetal/inventory/pkg/inventory"
)

func main() {
	is := inventory.NewInventorySvc()
	is.Inventorize()
}
