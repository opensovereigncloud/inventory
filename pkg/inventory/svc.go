package inventory

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/dev"
	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/flags"
	"github.com/onmetal/inventory/pkg/ioctl"
	"github.com/onmetal/inventory/pkg/pci"
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/proc"
	"github.com/onmetal/inventory/pkg/run"
	"github.com/onmetal/inventory/pkg/sys"
)

const (
	COKRetCode  = 0
	CErrRetCode = -1
)

type Svc struct {
	printer *printer.Svc

	dmiSvc     *dmi.Svc
	numaSvc    *sys.NumaSvc
	blockSvc   *sys.BlockSvc
	pciSvc     *sys.PCISvc
	cpuInfoSvc *proc.CPUInfoSvc
	memInfoSvc *proc.MemInfoSvc
	lldpSvc    *run.Svc
	nicSvc     *sys.NICSvc
	ipmiSvc    *ioctl.IPMISvc
	netlinkSvc *ioctl.NetlinkSvc
}

func NewSvc() (*Svc, int) {
	f := flags.NewFlags()

	p := printer.NewSvc(f.Verbose)

	pciIDs, err := pci.NewPCIIds()
	if err != nil {
		p.Err(errors.Wrapf(err, "unable to load PCI IDs"))
		return nil, CErrRetCode
	}

	rawDmiSvc := dmi.NewRawDMISvc(f.Root)
	dmiSvc := dmi.NewDMISvc(p, rawDmiSvc)

	cpuInfoSvc := proc.NewCPUInfoSvc(p, f.Root)
	memInfoSvc := proc.NewMemInfoSvc(p, f.Root)

	numaStatSvc := sys.NewNumaStatSvc(p)
	numaNodeSvc := sys.NewNumaNodeSvc(memInfoSvc, numaStatSvc)
	numaSvc := sys.NewNumaSvc(p, numaNodeSvc, f.Root)

	partitionTableSvc := dev.NewPartitionTableSvc(f.Root)
	blockDeviceStatSvc := sys.NewBlockDeviceStatSvc(p)
	blockDeviceSvc := sys.NewBlockDeviceSvc(p, partitionTableSvc, blockDeviceStatSvc)
	blockSvc := sys.NewBlockSvc(p, blockDeviceSvc, f.Root)

	pciDevSvc := sys.NewPCIDeviceSvc(p, pciIDs)
	pciBusSvc := sys.NewPCIBusSvc(p, pciDevSvc)
	pciSvc := sys.NewPCISvc(p, pciBusSvc, f.Root)

	lldpFrameInfoSvc := run.NewLLDPFrameInfoSvc(p)
	lldpSvc := run.NewLLDPSvc(p, lldpFrameInfoSvc, f.Root)

	nicDevSvc := sys.NewNICDeviceSvc(p)
	nicSvc := sys.NewNICSvc(p, nicDevSvc, f.Root)

	ipmiDevInfoSvc := ioctl.NewIPMIDeviceInfoSvc(p)
	ipmiSvc := ioctl.NewIPMISvc(p, ipmiDevInfoSvc, f.Root)

	nlSvc := ioctl.NewNetlinkSvc(p, f.Root)

	return &Svc{
		printer:    p,
		dmiSvc:     dmiSvc,
		numaSvc:    numaSvc,
		blockSvc:   blockSvc,
		pciSvc:     pciSvc,
		cpuInfoSvc: cpuInfoSvc,
		memInfoSvc: memInfoSvc,
		lldpSvc:    lldpSvc,
		nicSvc:     nicSvc,
		ipmiSvc:    ipmiSvc,
		netlinkSvc: nlSvc,
	}, 0
}

func (is *Svc) Inventorize() int {
	inv := &Inventory{}

	setters := []func(inventory *Inventory) error{
		is.setDMI,
		is.setCPUInfo,
		is.setMemInfo,
		is.setNumaNodes,
		is.setBlockDevices,
		is.setPCIBusDevices,
		is.setIPMIDevices,
		is.setNICs,
		is.setLLDPFrames,
		is.setNDPFrames,
	}

	for _, setter := range setters {
		err := setter(inv)
		if err != nil {
			is.printer.VErr(errors.Wrap(err, "unable to set value"))
		}
	}

	jsonBytes, err := json.Marshal(inv)
	if err != nil {
		is.printer.Err(errors.Wrap(err, "unable to marshal result to json"))
		return CErrRetCode
	}

	var prettifiedJsonBuf bytes.Buffer
	if err := json.Indent(&prettifiedJsonBuf, jsonBytes, "", "\t"); err != nil {
		is.printer.Err(errors.Wrap(err, "unable to indent json"))
		return CErrRetCode
	}

	is.printer.Out(prettifiedJsonBuf.String())

	return COKRetCode
}

func (is *Svc) setDMI(inv *Inventory) error {
	data, err := is.dmiSvc.GetDMIData()
	if err != nil {
		return errors.Wrap(err, "unable to get dmi data")
	}
	inv.DMI = data
	return nil
}

func (is *Svc) setCPUInfo(inv *Inventory) error {
	data, err := is.cpuInfoSvc.GetCPUInfo()
	if err != nil {
		return errors.Wrap(err, "unable to get proc data")
	}
	inv.CPUInfo = data
	return nil
}

func (is *Svc) setMemInfo(inv *Inventory) error {
	data, err := is.memInfoSvc.GetMemInfo()
	if err != nil {
		return errors.Wrap(err, "unable to get proc data")
	}
	inv.MemInfo = data
	return nil
}

func (is *Svc) setNumaNodes(inv *Inventory) error {
	data, err := is.numaSvc.GetNumaData()
	if err != nil {
		return errors.Wrap(err, "unable to get numa data")
	}
	inv.NumaNodes = data
	return nil
}

func (is *Svc) setBlockDevices(inv *Inventory) error {
	data, err := is.blockSvc.GetBlockData()
	if err != nil {
		return errors.Wrap(err, "unable to get block data")
	}
	inv.BlockDevices = data
	return nil
}

func (is *Svc) setPCIBusDevices(inv *Inventory) error {
	data, err := is.pciSvc.GetPCIData()
	if err != nil {
		return errors.Wrap(err, "unable to get pci data")
	}
	inv.PCIBusDevices = data
	return nil
}

func (is *Svc) setIPMIDevices(inv *Inventory) error {
	data, err := is.ipmiSvc.GetIPMIData()
	if err != nil {
		return errors.Wrap(err, "unable to get ipmi data")
	}
	inv.IPMIDevices = data
	return nil
}

func (is *Svc) setNICs(inv *Inventory) error {
	data, err := is.nicSvc.GetNICData()
	if err != nil {
		return errors.Wrap(err, "unable to get nic data")
	}
	inv.NICs = data
	return nil
}

func (is *Svc) setLLDPFrames(inv *Inventory) error {
	data, err := is.lldpSvc.GetLLDPData()
	if err != nil {
		return errors.Wrap(err, "unable to get lldp data")
	}
	inv.LLDPFrames = data
	return nil
}

func (is *Svc) setNDPFrames(inv *Inventory) error {
	data, err := is.netlinkSvc.GetIPv6NeighbourData()
	if err != nil {
		return errors.Wrap(err, "unable to get ndp data")
	}
	inv.NDPFrames = data
	return nil
}
