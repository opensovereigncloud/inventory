package numa

import (
	"github.com/onmetal/inventory/pkg/mem"
)

type Node struct {
	ID        int
	CPUs      []int
	Distances []int
	Memory    *mem.Info
	Stat      *Stat
}
