package lldp

import (
	"io/ioutil"
	"path"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CLLDPPath = "/run/systemd/netif/lldp"
)

type Svc struct {
	printer      *printer.Svc
	frameInfoSvc *FrameSvc
	lldpPath     string
}

func NewSvc(printer *printer.Svc, frameInfoSvc *FrameSvc, basePath string) *Svc {
	return &Svc{
		printer:      printer,
		frameInfoSvc: frameInfoSvc,
		lldpPath:     path.Join(basePath, CLLDPPath),
	}
}

func (s *Svc) GetData() ([]Frame, error) {
	frameFiles, err := ioutil.ReadDir(s.lldpPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of frame files")
	}

	frameInfos := make([]Frame, 0)

	// iterate over /run/systemd/netif/lldp/%i
	for _, frameFile := range frameFiles {
		fName := frameFile.Name()
		filePath := path.Join(s.lldpPath, fName)
		info, err := s.frameInfoSvc.GetFrame(fName, filePath)
		if err != nil {
			s.printer.VErr(errors.Errorf("unable to collect LLDP info for interface idx %s", fName))
			continue
		}
		frameInfos = append(frameInfos, *info)
	}

	return frameInfos, nil
}
