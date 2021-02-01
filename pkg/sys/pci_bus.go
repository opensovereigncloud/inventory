package sys

import (
	"io/ioutil"
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CPCIDeviceAddressPattern = "[[:xdigit:]]{4}:[[:xdigit:]]{2}:[[:xdigit:]]{2}.[[:xdigit:]]"
)

var CPCIDeviceAddressRegexp = regexp.MustCompile(CPCIDeviceAddressPattern)

type PCIBus struct {
	ID      string
	Devices []PCIDevice
}

func NewPCIBus(id string, devices []PCIDevice) *PCIBus {
	return &PCIBus{
		ID:      id,
		Devices: devices,
	}
}

type PCIBusSvc struct {
	printer      *printer.Svc
	pciDeviceSvc *PCIDeviceSvc
}

func NewPCIBusSvc(printer *printer.Svc, pciDevSvc *PCIDeviceSvc) *PCIBusSvc {
	return &PCIBusSvc{
		printer:      printer,
		pciDeviceSvc: pciDevSvc,
	}
}

func (s *PCIBusSvc) GetPCIBus(thePath string, id string) (*PCIBus, error) {
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

		device, err := s.pciDeviceSvc.GetPCIDevice(pciDevPath, name)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to collect info for device %s", name))
			continue
		}

		devices = append(devices, *device)
	}

	return NewPCIBus(id, devices), nil
}
