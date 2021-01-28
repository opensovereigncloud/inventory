package inventory

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/ioctl"
	"github.com/onmetal/inventory/pkg/proc"
	"github.com/onmetal/inventory/pkg/run"
	"github.com/onmetal/inventory/pkg/sys"
)

type Svc struct {
	dmiSvc     *dmi.Svc
	numaSvc    *sys.NumaSvc
	blockSvc   *sys.BlockSvc
	pciSvc     *sys.PCISvc
	procSvc    *proc.Svc
	lldpSvc    *run.Svc
	nicSvc     *sys.NICSvc
	ipmiSvc    *ioctl.IPMISvc
	netlinkSvc *ioctl.NetlinkSvc
}

func NewInventorySvc() *Svc {
	pciSvc, err := sys.NewPCISvc()
	if err != nil {
		panic(err)
	}

	return &Svc{
		dmiSvc:     dmi.NewDMISvc(),
		numaSvc:    sys.NewNumaSvc(),
		blockSvc:   sys.NewBlockSvc(),
		pciSvc:     pciSvc,
		procSvc:    proc.NewProcSvc(),
		lldpSvc:    run.NewLLDPSvc(),
		nicSvc:     sys.NewNICSvc(),
		ipmiSvc:    ioctl.NewIPMISvc(),
		netlinkSvc: ioctl.NewNetlinkSvc(),
	}
}

type Inventory struct {
	DMI           *dmi.DMI
	Proc          *proc.Proc
	NumaNodes     []sys.NumaNode
	BlockDevices  []sys.BlockDevice
	PCIBusDevices []sys.PCIBus
	IPMIDevices   []ioctl.IPMIDeviceInfo
	NICs          []sys.NIC
	LLDPFrames    []run.LLDPFrameInfo
	NDPFrames     []ioctl.IPv6Neighbour
}

func (is *Svc) Inventorize() {
	inv := &Inventory{}

	dmiData, err := is.dmiSvc.GetDMIData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.DMI = dmiData

	numaData, err := is.numaSvc.GetNumaData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.NumaNodes = numaData

	procData, err := is.procSvc.GetProcData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.Proc = procData

	blockData, err := is.blockSvc.GetBlockData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.BlockDevices = blockData

	pciData, err := is.pciSvc.GetPCIData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.PCIBusDevices = pciData

	lldpData, err := is.lldpSvc.GetLLDPData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.LLDPFrames = lldpData

	nicData, err := is.nicSvc.GetNICData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.NICs = nicData

	ipmiData, err := is.ipmiSvc.GetIPMIData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.IPMIDevices = ipmiData

	netData, err := is.netlinkSvc.GetIPv6NeighbourData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.NDPFrames = netData

	jsonBytes, err := json.Marshal(inv)
	if err != nil {
		fmt.Println(err)
		return
	}

	var prettifiedJsonBuf bytes.Buffer
	if err := json.Indent(&prettifiedJsonBuf, jsonBytes, "", "\t"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(prettifiedJsonBuf.String())
}
