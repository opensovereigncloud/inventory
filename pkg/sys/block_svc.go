package sys

import (
	"io/ioutil"
	"path"

	"github.com/pkg/errors"
)

const (
	CSysBlockBasePath = "/sys/block"
)

type BlockSvc struct{}

func NewBlockSvc() *BlockSvc {
	return &BlockSvc{}
}

func (bs *BlockSvc) GetBlockData() ([]BlockDevice, error) {
	blocks := make([]BlockDevice, 0)
	fileInfos, err := ioutil.ReadDir(CSysBlockBasePath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of block devices")
	}

	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		thePath := path.Join(CSysBlockBasePath, name)
		block, err := NewBlockDevice(thePath, name)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to collect block device data for %s", thePath)
		}

		blocks = append(blocks, *block)
	}

	return blocks, nil
}
