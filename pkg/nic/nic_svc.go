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

type NICSvc struct {
	printer    *printer.Svc
	nicDevSvc  *NICDeviceSvc
	nicDevPath string
}

func NewNICSvc(printer *printer.Svc, nicDevSvc *NICDeviceSvc, basePath string) *NICSvc {
	return &NICSvc{
		printer:    printer,
		nicDevSvc:  nicDevSvc,
		nicDevPath: path.Join(basePath, CNICDevicePath),
	}
}

func (s *NICSvc) GetNICData() ([]NIC, error) {
	nicFolders, err := ioutil.ReadDir(s.nicDevPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of nic folders")
	}

	var nics []NIC
	for _, nicFolder := range nicFolders {
		fName := nicFolder.Name()
		thePath := path.Join(s.nicDevPath, fName)
		nic, err := s.nicDevSvc.GetNIC(thePath, fName)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to collect NIC data"))
			continue
		}
		nics = append(nics, *nic)
	}

	return nics, nil
}
