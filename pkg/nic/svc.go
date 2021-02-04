package nic

import (
	"io/ioutil"
	"path"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CNICDevicePath = "/sys/class/net"
)

type Svc struct {
	printer    *printer.Svc
	nicDevSvc  *DeviceSvc
	nicDevPath string
}

func NewSvc(printer *printer.Svc, nicDevSvc *DeviceSvc, basePath string) *Svc {
	return &Svc{
		printer:    printer,
		nicDevSvc:  nicDevSvc,
		nicDevPath: path.Join(basePath, CNICDevicePath),
	}
}

func (s *Svc) GetData() ([]Device, error) {
	nicFolders, err := ioutil.ReadDir(s.nicDevPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of nic folders")
	}

	var nics []Device
	for _, nicFolder := range nicFolders {
		fName := nicFolder.Name()
		thePath := path.Join(s.nicDevPath, fName)
		nic, err := s.nicDevSvc.GetDevice(thePath, fName)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to collect Device data"))
			continue
		}
		nics = append(nics, *nic)
	}

	return nics, nil
}
