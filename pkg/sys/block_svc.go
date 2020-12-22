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

type Block struct {
	Blocks map[string]BlockDevice
}

func (bs *BlockSvc) GetBlockData() (*Block, error) {
	blocks := make(map[string]BlockDevice)
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

		blocks[name] = *block
	}

	return &Block{
		Blocks: blocks,
	}, nil
}
