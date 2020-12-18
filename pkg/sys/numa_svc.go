package sys

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type NumaSvc struct{}

func NewNumaSvc() *NumaSvc {
	return &NumaSvc{}
}

type Numa struct {
	Nodes map[int]NumaNode
}

const (
	CNodeDevicePath = "/sys/devices/system/node"

	CNumericNodeDeviceDirNamePattern = "node([0-9]+)"
)

var CNumericNodeDeviceDirNameRegexp = regexp.MustCompile(CNumericNodeDeviceDirNamePattern)

func (ns *NumaSvc) GetNumaData() (*Numa, error) {
	numaNodes := make(map[int]NumaNode)

	err := filepath.Walk(CNodeDevicePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "got error on directory traversal with path %s", path)
		}

		if !info.IsDir() {
			return nil
		}

		groups := CNumericNodeDeviceDirNameRegexp.FindStringSubmatch(info.Name())

		// String itself is always a first element in results
		// so we need at least 2 to get number from our group
		if len(groups) < 2 {
			return nil
		}

		nodeNumberString := groups[1]
		nodeNumber, err := strconv.Atoi(nodeNumberString)
		if err != nil {
			return errors.Wrapf(err, "unable to convert node number string %s to int", nodeNumberString)
		}

		node, err := NewNumaNode(path)
		if err != nil {
			return errors.Wrapf(err, "unable to collect  %s", path)
		}

		numaNodes[nodeNumber] = *node

		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "unable to walk through %s folder contents", CNodeDevicePath)
	}

	return &Numa{
		Nodes: numaNodes,
	}, nil
}
