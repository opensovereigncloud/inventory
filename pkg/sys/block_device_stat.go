package sys

import (
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

// https://www.kernel.org/doc/Documentation/block/stat.txt

const (
	CStatPath = "/stat"
)

type BlockDeviceStat struct {
	ReadIOs        uint64
	ReadMerges     uint64
	ReadSectors    uint64
	ReadTicks      uint64
	WriteIOs       uint64
	WriteMerges    uint64
	WriteSectors   uint64
	WriteTicks     uint64
	InFlight       uint64
	IOTicks        uint64
	TimeInQueue    uint64
	DiscardIOs     uint64
	DiscardMerges  uint64
	DiscardSectors uint64
	DiscardTicks   uint64
	FlushIOs       uint64
	FlushTicks     uint64
}

type BlockDeviceStatSvc struct {
	printer *printer.Svc
}

func NewBlockDeviceStatSvc(printer *printer.Svc) *BlockDeviceStatSvc {
	return &BlockDeviceStatSvc{
		printer: printer,
	}
}

func (s *BlockDeviceStatSvc) GetBlockDeviceStat(thePath string) (*BlockDeviceStat, error) {
	statPath := path.Join(thePath, CStatPath)
	contents, err := ioutil.ReadFile(statPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file %s", statPath)
	}

	stringContents := string(contents)
	trimmedStringContents := strings.TrimSpace(stringContents)

	fields := strings.Fields(trimmedStringContents)

	statVals := make([]uint64, len(fields))
	for i, field := range fields {
		val, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to convert to uint64 %s", field))
			val = 0
		}

		statVals[i] = val
	}

	stat := &BlockDeviceStat{}

	// linux kernel doc states that there are 11 fields
	// and underneath there is a table for 17
	// guess, we need to check this for the backward compatibility
	for i, val := range statVals {
		if err := stat.setByIndex(i, val); err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to set value %d on index %d", val, i))
		}
	}

	return stat, nil
}

func (s *BlockDeviceStat) setByIndex(idx int, val uint64) error {
	switch idx {
	case 0:
		s.ReadIOs = val
	case 1:
		s.ReadMerges = val
	case 2:
		s.ReadSectors = val
	case 3:
		s.ReadTicks = val
	case 4:
		s.WriteIOs = val
	case 5:
		s.WriteMerges = val
	case 6:
		s.WriteSectors = val
	case 7:
		s.WriteTicks = val
	case 8:
		s.InFlight = val
	case 9:
		s.IOTicks = val
	case 10:
		s.TimeInQueue = val
	case 11:
		s.DiscardIOs = val
	case 12:
		s.DiscardMerges = val
	case 13:
		s.DiscardSectors = val
	case 14:
		s.DiscardTicks = val
	case 15:
		s.FlushIOs = val
	case 16:
		s.FlushTicks = val
	default:
		return errors.Errorf("unexpected index %d", idx)
	}
	return nil
}
