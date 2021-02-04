package inventory

import (
	"github.com/onmetal/inventory/pkg/block"
	"github.com/onmetal/inventory/pkg/cpu"
	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/ipmi"
	"github.com/onmetal/inventory/pkg/lldp"
	"github.com/onmetal/inventory/pkg/mem"
	"github.com/onmetal/inventory/pkg/netlink"
	"github.com/onmetal/inventory/pkg/nic"
	"github.com/onmetal/inventory/pkg/numa"
	"github.com/onmetal/inventory/pkg/pci"
)

type Inventory struct {
	DMI           *dmi.DMI
	MemInfo       *mem.MemInfo
	CPUInfo       []cpu.Info
	NumaNodes     []numa.Node
	BlockDevices  []block.Device
	PCIBusDevices []pci.Bus
	IPMIDevices   []ipmi.IPMIDeviceInfo
	NICs          []nic.NIC
	LLDPFrames    []lldp.LLDPFrameInfo
	NDPFrames     []netlink.IPv6Neighbour
}
