package sys

import (
	"io/ioutil"
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/pci"
)

const (
	CDevicesPath = "/sys/devices"

	CPCIBusIDPattern = "pci(\\d{4}:\\d{2})"
)

var CPCIBusIDRegexp = regexp.MustCompile(CPCIBusIDPattern)

type PCISvc struct {
	ids *pci.IDs
}

func NewPCISvc() (*PCISvc, error) {
	ids, err := pci.NewPCIIds()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load PCI IDs")
	}

	return &PCISvc{
		ids: ids,
	}, nil
}

func (ps *PCISvc) GetPCIData() ([]PCIBus, error) {
	deviceFolders, err := ioutil.ReadDir(CDevicesPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of device folders")
	}

	var buses []PCIBus
	for _, deviceFolder := range deviceFolders {
		fName := deviceFolder.Name()

		groups := CPCIBusIDRegexp.FindStringSubmatch(fName)

		if len(groups) != 2 {
			continue
		}

		pciBusPath := path.Join(CDevicesPath, fName)
		bus, err := NewPCIBus(pciBusPath, groups[1], ps.ids)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to collect PCI bus %s data", groups[1])
		}

		buses = append(buses, *bus)
	}

	return buses, nil
}
