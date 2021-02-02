package sys

import (
	"io/ioutil"
	"path"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CSysBlockBasePath = "/sys/block"
)

type BlockSvc struct {
	printer      *printer.Svc
	devSvc       *BlockDeviceSvc
	sysBlockPath string
}

func NewBlockSvc(printer *printer.Svc, devSvc *BlockDeviceSvc, basePath string) *BlockSvc {
	return &BlockSvc{
		printer:      printer,
		devSvc:       devSvc,
		sysBlockPath: path.Join(basePath, CSysBlockBasePath),
	}
}

func (s *BlockSvc) GetBlockData() ([]BlockDevice, error) {
	blocks := make([]BlockDevice, 0)
	fileInfos, err := ioutil.ReadDir(s.sysBlockPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of block devices")
	}

	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		thePath := path.Join(s.sysBlockPath, name)
		block, err := s.devSvc.GetBlockDevice(thePath, name)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to collect block device data for %s", thePath))
		}

		blocks = append(blocks, *block)
	}

	return blocks, nil
}
