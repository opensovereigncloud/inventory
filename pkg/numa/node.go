package numa

import (
	"github.com/onmetal/inventory/pkg/mem"
)

type Node struct {
	ID       int
	CPUs     []int
	Distance int
	Memory   *mem.MemInfo
	Stat     *Stat
}
