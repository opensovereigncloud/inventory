package sys

import (
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/pci"
)

type PCISvc struct {
	ID *pci.IDs
}

func NewPCISvc() (*PCISvc, error) {
	ids, err := pci.NewPCIIds()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load PCI IDs")
	}

	return &PCISvc{
		ID: ids,
	}, nil
}
