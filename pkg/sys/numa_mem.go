package sys

import (
	"path"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/proc"
)

const (
	CNodeMemInfo = "/meminfo"
)

type NumaMemory proc.MemInfo

func NewNumaMemory(nodePath string) (*NumaMemory, error) {
	memInfoPath := path.Join(nodePath, CNodeMemInfo)
	mem, err := proc.NewMemInfoFromFile(memInfoPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get meminfo from %s", memInfoPath)
	}
	numaMem := NumaMemory(*mem)
	return &numaMem, err
}
