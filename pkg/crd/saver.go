package crd

import apiv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"

type SaverSvc interface {
	Save(inv *apiv1alpha1.Inventory) error
	Patch(name string, patch interface{}) error
}
