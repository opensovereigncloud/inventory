package gatherer

import (
	"github.com/onmetal/inventory/pkg/crd"
	"github.com/onmetal/inventory/pkg/inventory"
	"github.com/onmetal/inventory/pkg/printer"
	apiv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	"github.com/pkg/errors"
)

type NICUpdaterSvc struct {
	printer  *printer.Svc
	crdSvc   *crd.Svc
	gatherer *Svc
}

func NewNICUpdaterSvc() (*NICUpdaterSvc, int) {
	// TODO Make own flag set and service construction

	// TODO Extract coordination logic (flags, wiring, construction) into separate service (app)
	gatherer, _ := NewSvc()

	// TODO construct it yourself
	crdSvc := gatherer.crdSvc
	p := gatherer.printer

	return &NICUpdaterSvc{
		printer:  p,
		crdSvc:   crdSvc,
		gatherer: gatherer,
	}, 0
}

func (s *NICUpdaterSvc) Run() int {
	gatherSetters := []func(inventory *inventory.Inventory) error{
		s.gatherer.setDMI,
		s.gatherer.setNICs,
		s.gatherer.setLLDPFrames,
		s.gatherer.setNDPFrames,
		s.gatherer.setHost,
	}

	inv := s.gatherer.GatherInOrder(gatherSetters)

	// TODO separate build process from resource creation/update process
	buildSetters := []func(*apiv1alpha1.Inventory, *inventory.Inventory){
		s.crdSvc.SetSystem,
		s.crdSvc.SetNICs,
	}

	cr, err := s.crdSvc.BuildInOrder(inv, buildSetters)

	patch := struct {
		Spec struct {
			Nics []apiv1alpha1.NICSpec `json:"nics"`
		} `json:"spec"`
	}{}
	patch.Spec.Nics = cr.Spec.NICs

	err = s.crdSvc.Patch(cr.Name, patch)
	if err != nil {
		s.printer.Err(errors.Wrap(err, "unable to patch"))
		return CErrRetCode
	}

	return COKRetCode
}
