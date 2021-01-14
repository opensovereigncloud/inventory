package run

import (
	"io/ioutil"
	"path"

	"github.com/pkg/errors"
)

const (
	CLLDPPath = "/run/systemd/netif/lldp"
)

type LLDP struct {
	Frames []LLDPFrameInfo
}

type Svc struct{}

func NewLLDPSvc() *Svc {
	return &Svc{}
}

func (l *Svc) GetLLDPData() (*LLDP, error) {
	frameFiles, err := ioutil.ReadDir(CLLDPPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of frame files")
	}

	frameInfos := make([]LLDPFrameInfo, 0)

	// iterate over /run/systemd/netif/lldp/%i
	for _, frameFile := range frameFiles {
		fName := frameFile.Name()
		filePath := path.Join(CLLDPPath, fName)
		info, err := NewLLDPFrameInfo(fName, filePath)
		if err != nil {
			return nil, errors.Errorf("unable to collect LLDP info for interface idx %s", fName)
		}
		frameInfos = append(frameInfos, *info)
	}

	return &LLDP{
		Frames: frameInfos,
	}, nil
}
