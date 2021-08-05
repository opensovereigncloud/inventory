package benchmark

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/block"
	"github.com/onmetal/inventory/pkg/cpu"
	"github.com/onmetal/inventory/pkg/crd"
	"github.com/onmetal/inventory/pkg/flags"
	"github.com/onmetal/inventory/pkg/inventory"
	"github.com/onmetal/inventory/pkg/mem"
	"github.com/onmetal/inventory/pkg/mlc"
	"github.com/onmetal/inventory/pkg/numa"
	"github.com/onmetal/inventory/pkg/printer"
)

const (
	COKRetCode  = 0
	CErrRetCode = -1
)

type Svc struct {
	printer *printer.Svc

	crdSvc *crd.Svc

	numaSvc    *numa.Svc
	blockSvc   *block.Svc
	cpuInfoSvc *cpu.InfoSvc
	memInfoSvc *mem.InfoSvc
	mlcPerfSvc *mlcPerf.PerfSvc
}

func NewSvc() (*Svc, int) {
	f := flags.NewFlags()

	p := printer.NewSvc(f.Verbose)

	crdSvc, err := crd.NewSvc(f.Kubeconfig, f.KubeNamespace)
	if err != nil {
		p.Err(errors.Wrapf(err, "unable to create k8s resource svc"))
		return nil, CErrRetCode
	}

	cpuInfoSvc := cpu.NewInfoSvc(p, f.Root)
	memInfoSvc := mem.NewInfoSvc(p, f.Root)
	mlcPerfSvc := mlcPerf.NewPerfSvc(p, f.Root)

	numaStatSvc := numa.NewStatSvc(p)
	numaNodeSvc := numa.NewNodeSvc(memInfoSvc, numaStatSvc)
	numaSvc := numa.NewSvc(p, numaNodeSvc, f.Root)

	partitionTableSvc := block.NewPartitionTableSvc(f.Root)
	blockDeviceStatSvc := block.NewDeviceStatSvc(p)
	blockDeviceSvc := block.NewDeviceSvc(p, partitionTableSvc, blockDeviceStatSvc)
	blockSvc := block.NewSvc(p, blockDeviceSvc, f.Root)

	return &Svc{
		printer:    p,
		crdSvc:     crdSvc,
		numaSvc:    numaSvc,
		blockSvc:   blockSvc,
		cpuInfoSvc: cpuInfoSvc,
		mlcPerfSvc: mlcPerfSvc,
	}, 0
}

func (s *Svc) Gather() int {
	inv := &inventory.Inventory{}

	setters := []func(inventory *inventory.Inventory) error{
		s.setCPUInfo,
		s.setMlcPerf,
		s.setNumaNodes,
		s.setBlockDevices,
	}

	for _, setter := range setters {
		err := setter(inv)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to set value"))
		}
	}

	jsonBytes, err := json.Marshal(inv)
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to marshal result to json"))
	}

	var prettifiedJsonBuf bytes.Buffer
	if err := json.Indent(&prettifiedJsonBuf, jsonBytes, "", "\t"); err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to indent json"))
	}

	s.printer.VOut("Gathered data:")
	s.printer.VOut(prettifiedJsonBuf.String())

	if err := s.crdSvc.BuildAndSave(inv); err != nil {
		s.printer.Err(errors.Wrap(err, "unable to save inventory resource"))
		return CErrRetCode
	}

	return COKRetCode
}

func (s *Svc) setCPUInfo(inv *inventory.Inventory) error {
	data, err := s.cpuInfoSvc.GetInfo()
	if err != nil {
		return errors.Wrap(err, "unable to get proc data")
	}
	inv.CPUInfo = data
	return nil
}

func (s *Svc) setMlcPerf(inv *inventory.Inventory) error {
	data, err := s.mlcPerfSvc.GetInfo()
	if err != nil {
		return errors.Wrap(err, "unable to get mlc data")
	}
	inv.MlcPerf = data
	return nil
}

func (s *Svc) setNumaNodes(inv *inventory.Inventory) error {
	data, err := s.numaSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get numa data")
	}
	inv.NumaNodes = data
	return nil
}

func (s *Svc) setBlockDevices(inv *inventory.Inventory) error {
	data, err := s.blockSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get block data")
	}
	inv.BlockDevices = data
	return nil
}
