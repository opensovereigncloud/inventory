package inventory

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/proc"
	"github.com/onmetal/inventory/pkg/run"
	"github.com/onmetal/inventory/pkg/sys"
)

type Svc struct {
	dmiSvc   *dmi.Svc
	numaSvc  *sys.NumaSvc
	blockSvc *sys.BlockSvc
	pciSvc   *sys.PCISvc
	procSvc  *proc.Svc
	lldpSvc  *run.Svc
	nicSvc   *sys.NICSvc
}

func NewInventorySvc() *Svc {
	pciSvc, err := sys.NewPCISvc()
	if err != nil {
		panic(err)
	}

	return &Svc{
		dmiSvc:   dmi.NewDMISvc(),
		numaSvc:  sys.NewNumaSvc(),
		blockSvc: sys.NewBlockSvc(),
		pciSvc:   pciSvc,
		procSvc:  proc.NewProcSvc(),
		lldpSvc:  run.NewLLDPSvc(),
		nicSvc:   sys.NewNICSvc(),
	}
}

type Inventory struct {
	DMI     *dmi.DMI
	Numa    *sys.Numa
	Block   *sys.Block
	Proc    *proc.Proc
	PCI     *sys.PCI
	LLDP    *run.LLDP
	Network *sys.Network
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
	inv.Numa = numaData

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
	inv.Block = blockData

	pciData, err := is.pciSvc.GetPCIData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.PCI = pciData

	lldpData, err := is.lldpSvc.GetLLDPData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.LLDP = lldpData

	nicData, err := is.nicSvc.GetNICData()
	if err != nil {
		fmt.Println(err)
		return
	}
	inv.Network = nicData

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
