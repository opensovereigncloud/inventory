package mlcPerf

import (
	"github.com/pkg/errors"
)

const (
	CMemPerfLocalMemBWKey       = "LocalMemBW"
	CMemPerfRemoteMemBWKey      = "RemoteMemBW"
	CMemPerfLocalMemLatencyKey  = "LocalMemLatency"
	CMemPerfRemoteMemLatencyKey = "RemoteMemLatency"
)

type Perf struct {
	LocalMemBW       float64
	RemoteMemBW      float64
	LocalMemLatency  float64
	RemoteMemLatency float64
}

func (memperf *Perf) setField(key string, val float64) error {
	switch key {
	case CMemPerfLocalMemBWKey:
		memperf.LocalMemBW = val
	case CMemPerfRemoteMemBWKey:
		memperf.RemoteMemBW = val
	case CMemPerfLocalMemLatencyKey:
		memperf.LocalMemLatency = val
	case CMemPerfRemoteMemLatencyKey:
		memperf.RemoteMemLatency = val
	default:
		return errors.Errorf("unknown key %s from meminfo", key)
	}
	return nil
}
