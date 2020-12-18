package inventory

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/proc"
	"github.com/onmetal/inventory/pkg/sys"
)

type Svc struct {
	dmiSvc  *dmi.Svc
	numaSvc *sys.NumaSvc
	procSvc *proc.Svc
}

func NewInventorySvc() *Svc {
	return &Svc{
		dmiSvc:  dmi.NewDMISvc(),
		numaSvc: sys.NewNumaSvc(),
		procSvc: proc.NewProcSvc(),
	}
}

type Inventory struct {
	DMI  *dmi.DMI
	Numa *sys.Numa
	Proc *proc.Proc
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
