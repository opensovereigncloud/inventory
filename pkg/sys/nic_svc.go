package sys

import (
	"io/ioutil"
	"path"

	"github.com/pkg/errors"
)

const (
	CNICDevicePath = "/sys/class/net"
)

type Network struct {
	NICs []NIC
}

type NICSvc struct{}

func NewNICSvc() *NICSvc {
	return &NICSvc{}
}

func (ns *NICSvc) GetNICData() (*Network, error) {
	nicFolders, err := ioutil.ReadDir(CNICDevicePath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of nic folders")
	}

	var nics []NIC
	for _, nicFolder := range nicFolders {
		fName := nicFolder.Name()
		thePath := path.Join(CNICDevicePath, fName)
		nic, err := NewNIC(thePath, fName)
		if err != nil {
			return nil, errors.Wrap(err, "unable to collect NIC data")
		}
		nics = append(nics, *nic)
	}

	return &Network{
		NICs: nics,
	}, nil
}
