package virt

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jeek120/cpuid"
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/cpu"
	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/file"
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

	CProcXenPath                        = "/proc/xen"
	CSysHypervisorTypePaht              = "/sys/hypervisor/type"
	CDeviceTreePath                     = "/proc/device-tree"
	CDeviceTreeHypervisorCompatiblePath = "/proc/device-tree/hypervisor/compatible"
	CDeviceTreeIBMPartitionNamePath     = "/proc/device-tree/ibm,partition-name"
	CDeviceTreeHMCManagedPath           = "/proc/device-tree/hmc-managed?"
	CDeviceTreeQEMUPath                 = "/proc/device-tree/chosen/qemu,graphic-width"
	CProcSysInfo                        = "/proc/sysinfo"
)

type Type string

type Virtualization struct {
	Type Type
}

type Svc struct {
	dmiSvc     *dmi.Svc
	cpuInfoSvc *cpu.InfoSvc

	procXenPath                        string
	sysHypervisorTypePath              string
	deviceTreePath                     string
	deviceTreeHypervisorCompatiblePath string
	deviceTreeIBMPartitionNamePath     string
	deviceTreeHMCManagedPath           string
	deviceTreeQEMUPath                 string
	procSysInfoPath                    string
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

var CCPUIDVMStrings = map[string]Type{
	"XenVMMXenVMM": CTypeVMXen,
	"KVMKVMKVM":    CTypeVMKVM,
	"TCGTCGTCGTCG": CTypeVMQEMU,
	"VMwareVMware": CTypeVMVMware,
	"Microsoft Hv": CTypeVMMicrosoft,
	"bhyve bhyve ": CTypeVMBhyve,
	"QNXQVMBSQG":   CTypeVMQNX,
	"ACRNACRNACRN": CTypeVMACRN,
}

func (s *Svc) GetData() {
	theType := s.checkVMWithDMI()

	if theType == CTypeVMXen || theType == CTypeVMOracle {

	}
}

func (s *Svc) checkVMWithDMI() Type {
	dmiData, _ := s.dmiSvc.GetData()

	vendorLocators := []string{
		dmiData.SystemInformation.ProductName,
		dmiData.SystemInformation.Manufacturer,
		dmiData.BIOSInformation.Vendor,
	}
	for _, board := range dmiData.BoardInformation {
		vendorLocators = append(vendorLocators, board.Manufacturer)
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

func (s *Svc) checkVMWithCPUInfo() Type {
	infos, _ := s.cpuInfoSvc.GetInfo()

	for _, info := range infos {
		if strings.HasPrefix(info.VendorID, "User Mode Linux") {
			return CTypeVMUML
		}
	}

	return CTypeNone
}

func (s *Svc) checkVMWithCPUID() (Type, error) {
	ids := [4]uint32{}
	cpuid.Cpuid(&ids, 1)

	hypervisor := ids[2] & 0x80000000

	if hypervisor == 0 {
		return CTypeNone, nil
	}

	ids = [4]uint32{}
	cpuid.Cpuid(&ids, 0x40000000)

	buf := new(bytes.Buffer)
	for i := 1; i < 4; i++ {
		if err := binary.Write(buf, binary.LittleEndian, ids[i]); err != nil {
			return "", errors.Wrap(err, "unable to write to buffer")
		}
	}

	text := string(buf.Bytes())

	for k, v := range CCPUIDVMStrings {
		if strings.HasPrefix(text, k) {
			return v, nil
		}
	}

	return CTypeVMOther, nil
}

func (s *Svc) checkVMWithXen() Type {
	if _, err := os.Stat(s.procXenPath); os.IsNotExist(err) {
		return CTypeNone
	}

	return CTypeVMXen
}

func (s *Svc) checkVMWithSysHypervisor() Type {
	if _, err := os.Stat(s.sysHypervisorTypePath); os.IsNotExist(err) {
		return CTypeNone
	}

	str, err := file.ToString(s.sysHypervisorTypePath)
	if err != nil {
		return errors.Wrapf(err, "unable to read %s into string", s.sysHypervisorTypePath)
	}

	if str == "xen" {
		return CTypeVMXen
	}

	return CTypeVMOther
}

func (s *Svc) checkVMWithDeviceTree() (Type, error) {
	_, err := os.Stat(s.deviceTreeHypervisorCompatiblePath)
	if !os.IsNotExist(err) {
		return "", errors.Wrapf(err, "unable to open file %s", s.deviceTreeHypervisorCompatiblePath)
	} else {
		access := 0

		if _, err := os.Stat(s.deviceTreeIBMPartitionNamePath); err != nil {
			access++
		}
		if _, err := os.Stat(s.deviceTreeHMCManagedPath); err != nil {
			access++
		}
		if _, err := os.Stat(s.deviceTreeQEMUPath); err != nil {
			access++
		}

		if access == 0 {
			return CTypeVMPowerVM, nil
		}

		files, err := ioutil.ReadDir(s.deviceTreePath)
		if os.IsNotExist(err) {
			return CTypeNone, nil
		}
		if err != nil {
			return "", errors.Wrapf(err, "unable to open directory %s", s.deviceTreePath)
		}

		for _, f := range files {
			if strings.Contains(f.Name(), "fw-cfg") {
				return CTypeVMQEMU, nil
			}
		}
	}

	str, err := file.ToString(s.deviceTreeHypervisorCompatiblePath)
	if err != nil {
		return "", errors.Wrapf(err, "unable to open file %s", s.deviceTreeHypervisorCompatiblePath)
	}

	if str == "linux,kvm" {
		return CTypeVMKVM, nil
	}
	if strings.Contains(str, "xen") {
		return CTypeVMXen, nil
	}
	if strings.Contains(str, "vmware") {
		return CTypeVMXen, nil
	}

	return CTypeVMOther, nil
}

func (s *Svc) checkVMWithZVM() (Type, error) {
	sysInfoData, err := ioutil.ReadFile(s.procSysInfoPath)
	if os.IsNotExist(err) {
		return CTypeNone, nil
	}
	if err != nil {
		return "", errors.Wrapf(err, "unable to read cpu info from %s", s.procSysInfoPath)
	}

	bufReader := bytes.NewReader(sysInfoData)
	scanner := bufio.NewScanner(bufReader)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "VM00 Control Program") {
			continue
		}
		if strings.Contains(line, "z/VM") {
			return CTypeVMZVM, nil
		} else {
			return CTypeVMKVM, nil
		}
	}

	return CTypeNone, nil
}
