package ioctl

import (
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

const (
	CDevPath        = "/dev"
	CIPMIDevPattern = "ipmi(\\d+)"
)

var CIPMIDevRegexp = regexp.MustCompile(CIPMIDevPattern)

type IPMISvc struct{}

func NewIPMISvc() *IPMISvc {
	return &IPMISvc{}
}

func (s *IPMISvc) GetIPMIData() ([]IPMIDeviceInfo, error) {
	devFolderContents, err := ioutil.ReadDir(CDevPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read contents of %s", CDevPath)
	}

	infos := make([]IPMIDeviceInfo, 0)
	for _, dev := range devFolderContents {
		devName := dev.Name()

		groups := CIPMIDevRegexp.FindStringSubmatch(devName)

		if len(groups) != 2 {
			continue
		}

		numStr := groups[1]
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, errors.Wrapf(err, "unabale to convert %s to int", numStr)
		}

		info, err := NewIPMIDeviceInfo(num)
		if err != nil {
			return nil, errors.Wrapf(err, "unabale to convert %s to int", numStr)
		}

		infos = append(infos, *info)
	}

	return infos, nil
}
