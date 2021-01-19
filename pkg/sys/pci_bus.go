package sys

import (
	"io/ioutil"
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/pci"
)

const (
	CPCIDeviceAddressPattern = "[[:xdigit:]]{4}:[[:xdigit:]]{2}:[[:xdigit:]]{2}.[[:xdigit:]]"
)

var CPCIDeviceAddressRegexp = regexp.MustCompile(CPCIDeviceAddressPattern)

type PCIBus struct {
	ID      string
	Devices []PCIDevice
}

func NewPCIBus(thePath string, id string, ids *pci.IDs) (*PCIBus, error) {
	pciDevFolders, err := ioutil.ReadDir(thePath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of block devices")
	}

	devices := make([]PCIDevice, 0)
	for _, pciDevFolder := range pciDevFolders {
		name := pciDevFolder.Name()

		if !pciDevFolder.IsDir() {
			continue
		}

		match := CPCIDeviceAddressRegexp.MatchString(name)
		if !match {
			continue
		}

		pciDevPath := path.Join(thePath, name)

		device, err := NewPCIDevice(pciDevPath, name, ids)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to collect info for device %s", name)
		}

		devices = append(devices, *device)
	}

	return &PCIBus{
		ID:      id,
		Devices: devices,
	}, nil
}
