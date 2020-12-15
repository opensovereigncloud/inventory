package inventory

import (
	"encoding/json"
	"fmt"

	"github.com/onmetal/inventory/pkg/dmi"
)

type Svc struct {
}

func NewInventorySvc() *Svc {
	return &Svc{}
}

func (is *Svc) Inventorize() {
	dmiSvc := dmi.NewDMISvc()

	dmiData, err := dmiSvc.GetDMIData()
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonBytes, err := json.Marshal(dmiData)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonBytes))
}
