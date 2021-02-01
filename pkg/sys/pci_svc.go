package sys

import (
	"io/ioutil"
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CDevicesPath = "/sys/devices"

	CPCIBusIDPattern = "pci(\\d{4}:\\d{2})"
)

var CPCIBusIDRegexp = regexp.MustCompile(CPCIBusIDPattern)

type PCISvc struct {
	printer     *printer.Svc
	pciBusSvc   *PCIBusSvc
	devicesPath string
}

func NewPCISvc(printer *printer.Svc, pciBusSvc *PCIBusSvc, basePath string) *PCISvc {
	return &PCISvc{
		printer:     printer,
		pciBusSvc:   pciBusSvc,
		devicesPath: path.Join(basePath, CDevicesPath),
	}
}

func (s *PCISvc) GetPCIData() ([]PCIBus, error) {
	deviceFolders, err := ioutil.ReadDir(s.devicesPath)
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

		pciBusPath := path.Join(s.devicesPath, fName)
		bus, err := s.pciBusSvc.GetPCIBus(pciBusPath, groups[1])
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to collect PCI bus %s data", groups[1]))
			continue
		}

		buses = append(buses, *bus)
	}

	return buses, nil
}
