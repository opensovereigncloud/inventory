package sys

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CNodeStat = "/numastat"

	CNodeNumaHitKey       = "numa_hit"
	CNodeNumaMissKey      = "numa_miss"
	CNodeNumaForeignKey   = "numa_foreign"
	CNodeInterleaveHitKey = "interleave_hit"
	CNodeLocalNodeKey     = "local_node"
	CNodeOtherNodeKey     = "other_node"

	CNodeStatLinePattern = "^(\\w+)\\s+(\\d+)$"
)

var CNodeStatLineRegexp = regexp.MustCompile(CNodeStatLinePattern)

type NumaStat struct {
	NumaHit       uint64
	NumaMiss      uint64
	NumaForeign   uint64
	InterleaveHit uint64
	LocalNode     uint64
	OtherNode     uint64
}

type NumaStatSvc struct {
	printer *printer.Svc
}

func NewNumaStatSvc(printer *printer.Svc) *NumaStatSvc {
	return &NumaStatSvc{
		printer: printer,
	}
}

func (s *NumaStatSvc) GetNumaStat(thePath string) (*NumaStat, error) {
	statPath := path.Join(thePath, CNodeStat)
	statData, err := ioutil.ReadFile(statPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read stat file from path %s", statPath)
	}

	stat := &NumaStat{}

	bufReader := bytes.NewReader(statData)
	scanner := bufio.NewScanner(bufReader)
	for scanner.Scan() {
		line := scanner.Text()

		groups := CNodeStatLineRegexp.FindStringSubmatch(line)

		// groups [0] string [1] key [2] value
		if len(groups) < 3 {
			continue
		}

		key := groups[1]
		valString := groups[2]

		val, err := strconv.ParseUint(valString, 10, 64)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to parse %s:%s into uint64", key, valString))
			continue
		}

		err = stat.setField(key, val)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to set %s:%d", key, val))
		}
	}

	return stat, nil
}

func (stat *NumaStat) setField(key string, val uint64) error {
	switch key {
	case CNodeNumaHitKey:
		stat.NumaHit = val
	case CNodeNumaMissKey:
		stat.NumaMiss = val
	case CNodeNumaForeignKey:
		stat.NumaForeign = val
	case CNodeInterleaveHitKey:
		stat.InterleaveHit = val
	case CNodeLocalNodeKey:
		stat.LocalNode = val
	case CNodeOtherNodeKey:
		stat.OtherNode = val
	default:
		return errors.Errorf("unknown key %s from meminfo", key)
	}
	return nil
}
