package virt

import (
	"strings"

	"github.com/onmetal/inventory/pkg/dmi"
)

const (
	CTypeNone        = "none"
	CTypeVMKVM       = "kvm"
	CTypeVMQEMU      = "qemu"
	CTypeVMBochs     = "bochs"
	CTypeVMXen       = "xen"
	CTypeVMUML       = "uml"
	CTypeVMVMware    = "vmware"
	CTypeVMOracle    = "oracle"
	CTypeVMMicrosoft = "microsoft"
	CTypeVMZVM       = "zvm"
	CTypeVMParallels = "parallels"
	CTypeVMBhyve     = "bhyve"
	CTypeVMQNX       = "qnx"
	CTypeVMACRN      = "acrn"
	CTypeVMPowerVM   = "powervm"
	CTypeVMOther     = "other"
)

type Type string

type Virtualization struct {
	Type Type
}

type Svc struct {
	dmiSvc *dmi.Svc
}

func NewSvc() *Svc {
	return &Svc{}
}

var CDMIVendorPrefixes = map[string]Type{
	"KVM":                CTypeVMKVM,
	"QEMU":               CTypeVMQEMU,
	"VMware":             CTypeVMVMware,
	"VMW":                CTypeVMVMware,
	"innotek GmbH":       CTypeVMOracle,
	"Oracle Corporation": CTypeVMOracle,
	"Xen":                CTypeVMXen,
	"Bochs":              CTypeVMBochs,
	"Parallels":          CTypeVMParallels,
	"BHYVE":              CTypeVMBhyve,
}

func (s *Svc) GetData() {
	theType := s.checkVMWIthDMI()

	if theType == CTypeVMXen || theType == CTypeVMOracle {

	}
}

func (s *Svc) checkVMWIthDMI() Type {
	dmiData, _ := s.dmiSvc.GetData()

	vendorLocators := []string{
		dmiData.SystemInformation.ProductName,
		dmiData.SystemInformation.Manufacturer,
		// TODO: inset board vendor here
		dmiData.BIOSInformation.Vendor,
	}

	for _, vendorLocator := range vendorLocators {
		for k, v := range CDMIVendorPrefixes {
			if strings.HasPrefix(vendorLocator, k) {
				return v
			}
		}
	}

	return CTypeNone
}
