package inventory

import (
	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/ioctl"
	"github.com/onmetal/inventory/pkg/proc"
	"github.com/onmetal/inventory/pkg/run"
	"github.com/onmetal/inventory/pkg/sys"
)

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
