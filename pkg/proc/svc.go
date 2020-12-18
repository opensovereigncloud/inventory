package proc

import "github.com/pkg/errors"

type Proc struct {
	MemInfo *MemInfo
	CPUInfo []CPUInfo
}

type Svc struct{}

func NewProcSvc() *Svc {
	return &Svc{}
}

func (ps *Svc) GetProcData() (*Proc, error) {
	mem, err := NewMemInfo()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get meminfo")
	}

	cpu, err := NewCPUInfo()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get cpuinfo")
	}

	return &Proc{
		MemInfo: mem,
		CPUInfo: cpu,
	}, nil
}
