package nic

import (
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/file"
	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CNICDevicePCIAddressPath = "/device"

	CNICDeviceAddressAddressAssignTypePath = "/addr_assign_type"
	CNICDeviceAddressPath                  = "/address"
	CNICDeviceAddressLengthPath            = "/addr_len"
	CNICDeviceBroadcastPath                = "/broadcast"
	CNICDeviceCarrierPath                  = "/carrier"
	CNICDeviceCarrierChangesPath           = "/carrier_changes"
	CNICDeviceCarrierDownCountPath         = "/carrier_down_count"
	CNICDeviceCarrierUpCountPath           = "/carrier_up_count"
	CNICDeviceDevIDPath                    = "/dev_id"
	CNICDeviceDevPortPath                  = "/dev_port"
	CNICDeviceDormantPath                  = "/dormant"
	CNICDeviceDuplexPath                   = "/duplex"
	CNICDeviceFlagsPath                    = "/flags"
	CNICDeviceInterfaceAliasPath           = "/ifalias"
	CNICDeviceInterfaceIndexPath           = "/ifindex"
	CNICDeviceInterfaceLinkPath            = "/iflink"
	CNICDeviceLinkModePath                 = "/link_mode"
	CNICDeviceMTUPath                      = "/mtu"
	CNICDeviceNameAssignTypePath           = "/name_assign_type"
	CNICDeviceNetDevGroupPath              = "/netdev_group"
	CNICDeviceOperationalStatePath         = "/operstate"
	CNICDevicePhysicalPortIDPath           = "/phys_port_id"
	CNICDevicePhysicalPortNamePath         = "/phys_port_name"
	CNICDevicePhysicalSwitchIDPath         = "/phys_switch_id"
	CNICDeviceSpeedPath                    = "/speed"
	CNICDeviceTestingPath                  = "/testing"
	CNICDeviceTransmitQueueLengthPath      = "/tx_queue_len"
	CNICDeviceTypePath                     = "/type"

	CPermanentAddressAssignType               = "permanent address"
	CRandomlyGeneratedAddressAssignType       = "randomly generated"
	CStolenFromAnotherDeviceAddressAssignType = "stolen from another device"
	CSetUsingDevAddressAssignType             = "set using dev_set_mac_address"

	CDefaultLinkMode = "default link mode"
	CDormantLinkMode = "dormant link mode"

	CUnpredictableKernelNameAssignType = "enumerated by the kernel, possibly in an unpredictable way"
	CPredictableKernelNameAssignType   = "predictably named by the kernel"
	CUserspaceNameAssignType           = "named by userspace"
	CRenamedNameAssignType             = "renamed"
)

type NICAddressAssignType string

type NICLinkMode string

type NICNameAssignType string

var CAddressAssignTypes = []NICAddressAssignType{
	CPermanentAddressAssignType,
	CRandomlyGeneratedAddressAssignType,
	CStolenFromAnotherDeviceAddressAssignType,
	CSetUsingDevAddressAssignType,
}

var CLinkModes = []NICLinkMode{
	CDefaultLinkMode,
	CDormantLinkMode,
}

var CNameAssignTypes = []NICNameAssignType{
	CUnpredictableKernelNameAssignType,
	CPredictableKernelNameAssignType,
	CUserspaceNameAssignType,
	CRenamedNameAssignType,
}

type NIC struct {
	Name       string
	PCIAddress string

	AddressAssignType   NICAddressAssignType
	Address             string
	AddressLength       uint8
	Broadcast           string
	Carrier             bool
	CarrierChanges      uint32
	CarrierDownCount    uint32
	CarrierUpCount      uint32
	DevID               string
	DevPort             uint8
	Dormant             bool
	Duplex              string
	Flags               *NICFlags
	InterfaceAlias      string
	InterfaceIndex      uint32
	InterfaceLink       uint32
	LinkMode            NICLinkMode
	MTU                 uint16
	NetDevGroup         int
	OperationalState    string
	PhysicalPortID      string
	PhysicalPortName    string
	PhysicalSwitchID    string
	Speed               uint32
	Testing             bool
	TransmitQueueLength uint32
	Type                CNICType
}

type NICDeviceSvc struct {
	printer *printer.Svc
}

func NewNICDeviceSvc(printer *printer.Svc) *NICDeviceSvc {
	return &NICDeviceSvc{
		printer: printer,
	}
}

func (s *NICDeviceSvc) GetNIC(thePath string, name string) (*NIC, error) {
	nic := &NIC{
		Name: name,
	}

	defs := []func(string) error{
		nic.defPCIAddress,
		nic.defAddressAssignType,
		nic.defAddress,
		nic.defAddressLength,
		nic.defBroadcast,
		nic.defCarrier,
		nic.defCarrierChanges,
		nic.defCarrierDownCount,
		nic.defCarrierUpCount,
		nic.defDevID,
		nic.defDevPort,
		nic.defDormant,
		nic.defDuplex,
		nic.defFlags,
		nic.defInterfaceAlias,
		nic.defInterfaceIndex,
		nic.defInterfaceLink,
		nic.defLinkMode,
		nic.defMTU,
		nic.defNetDevGroup,
		nic.defOperationalState,
		nic.defPhysicalPortID,
		nic.defPhysicalPortName,
		nic.defPhysicalSwitchID,
		nic.defSpeed,
		nic.defTesting,
		nic.defTransmitQueueLength,
		nic.defType,
	}

	for _, def := range defs {
		err := def(thePath)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to set NIC property"))
		}
	}

	return nic, nil
}

func (n *NIC) defPCIAddress(thePath string) error {
	filePath := path.Join(thePath, CNICDevicePCIAddressPath)

	linkPath, err := filepath.EvalSymlinks(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable resolve symlink with path %s", filePath)
	}
	fileInfo, err := os.Stat(linkPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get stat for path %s", filePath)
	}

	n.PCIAddress = fileInfo.Name()

	return nil
}

func (n *NIC) defAddressAssignType(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceAddressAddressAssignTypePath)
	fileVal, err := file.ToInt(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get address assign type from %s", filePath)
	}

	if fileVal < 0 || fileVal >= len(CAddressAssignTypes) {
		return errors.Errorf("unexpected value %d for address assign type", fileVal)
	}
	n.AddressAssignType = CAddressAssignTypes[fileVal]

	return nil
}

func (n *NIC) defAddress(thePath string) error {
	addressPath := path.Join(thePath, CNICDeviceAddressPath)
	addressString, err := file.ToString(addressPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get address string from %s", addressPath)
	}

	n.Address = addressString

	return nil
}

func (n *NIC) defAddressLength(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceAddressLengthPath)
	fileVal, err := file.ToUint8(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get address length from %s", filePath)
	}

	n.AddressLength = fileVal

	return nil
}

func (n *NIC) defBroadcast(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceBroadcastPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get broadcast address from %s", filePath)
	}

	n.Broadcast = fileVal

	return nil
}

func (n *NIC) defCarrier(thePath string) error {
	carrierPath := path.Join(thePath, CNICDeviceCarrierPath)
	carrier, err := file.ToBool(carrierPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get carrier status from %s", carrierPath)
	}

	n.Carrier = carrier

	return nil
}

func (n *NIC) defCarrierChanges(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceCarrierChangesPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get carrier changes address from %s", filePath)
	}

	n.CarrierChanges = fileVal

	return nil
}

func (n *NIC) defCarrierDownCount(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceCarrierDownCountPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get carrier down count from %s", filePath)
	}

	n.CarrierChanges = fileVal

	return nil
}

func (n *NIC) defCarrierUpCount(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceCarrierUpCountPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get carrier up count from %s", filePath)
	}

	n.CarrierChanges = fileVal

	return nil
}

func (n *NIC) defDevID(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceDevIDPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get dev ID from %s", filePath)
	}

	n.DevID = fileVal

	return nil
}

func (n *NIC) defDevPort(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceDevPortPath)
	fileVal, err := file.ToUint8(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get dev port from %s", filePath)
	}

	n.DevPort = fileVal

	return nil
}

func (n *NIC) defDormant(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceDormantPath)
	fileVal, err := file.ToBool(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get dormant value from %s", filePath)
	}

	n.Dormant = fileVal

	return nil
}

func (n *NIC) defDuplex(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceDuplexPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get duplex value from %s", filePath)
	}

	n.Duplex = fileVal

	return nil
}

func (n *NIC) defFlags(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceFlagsPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get flags value from %s", filePath)
	}

	val, err := strconv.ParseUint(fileVal[2:], 16, 32)
	if err != nil {
		return errors.Wrapf(err, "unable to parse uint from from %s", fileVal)
	}

	n.Flags = NewNICFlags(uint32(val))

	return nil
}

func (n *NIC) defInterfaceAlias(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceInterfaceAliasPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get interface alias value from %s", filePath)
	}

	n.InterfaceAlias = fileVal

	return nil
}

func (n *NIC) defInterfaceIndex(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceInterfaceIndexPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get interface index value from %s", filePath)
	}

	n.InterfaceIndex = fileVal

	return nil
}

func (n *NIC) defInterfaceLink(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceInterfaceLinkPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get interface link value from %s", filePath)
	}

	n.InterfaceLink = fileVal

	return nil
}

func (n *NIC) defLinkMode(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceLinkModePath)
	fileVal, err := file.ToInt(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get link mode from %s", filePath)
	}

	if fileVal < 0 || fileVal >= len(CLinkModes) {
		return errors.Errorf("unexpected value %d for link mode", fileVal)
	}
	n.LinkMode = CLinkModes[fileVal]

	return nil
}

func (n *NIC) defMTU(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceMTUPath)
	fileVal, err := file.ToUint16(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get MTU value from %s", filePath)
	}

	n.MTU = fileVal

	return nil
}

func (n *NIC) defNameAssignType(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceNameAssignTypePath)
	fileVal, err := file.ToInt(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get name assign type from %s", filePath)
	}

	fileVal = fileVal - 1

	if fileVal < 0 || fileVal >= len(CNameAssignTypes) {
		return errors.Errorf("unexpected value %d for name assign type", fileVal)
	}
	n.LinkMode = CLinkModes[fileVal]

	return nil
}

func (n *NIC) defNetDevGroup(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceNetDevGroupPath)
	fileVal, err := file.ToInt(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get network device group value from %s", filePath)
	}

	n.NetDevGroup = fileVal

	return nil
}

func (n *NIC) defOperationalState(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceOperationalStatePath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get operational state value from %s", filePath)
	}

	n.OperationalState = fileVal

	return nil
}

func (n *NIC) defPhysicalPortID(thePath string) error {
	filePath := path.Join(thePath, CNICDevicePhysicalPortIDPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get physical port ID value from %s", filePath)
	}

	n.PhysicalPortID = fileVal

	return nil
}

func (n *NIC) defPhysicalPortName(thePath string) error {
	filePath := path.Join(thePath, CNICDevicePhysicalPortNamePath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get physical port name value from %s", filePath)
	}

	n.PhysicalPortName = fileVal

	return nil
}

func (n *NIC) defPhysicalSwitchID(thePath string) error {
	filePath := path.Join(thePath, CNICDevicePhysicalSwitchIDPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get physical switch ID value from %s", filePath)
	}

	n.PhysicalSwitchID = fileVal

	return nil
}

func (n *NIC) defSpeed(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceSpeedPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get speed value from %s", filePath)
	}

	n.Speed = fileVal

	return nil
}

func (n *NIC) defTesting(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceTestingPath)
	fileVal, err := file.ToBool(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get testing value from %s", filePath)
	}

	n.Testing = fileVal

	return nil
}

func (n *NIC) defTransmitQueueLength(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceTransmitQueueLengthPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get transmit queue length value from %s", filePath)
	}

	n.TransmitQueueLength = fileVal

	return nil
}

func (n *NIC) defType(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceTypePath)
	fileVal, err := file.ToUint16(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get type value from %s", filePath)
	}

	theType, ok := CNICTypes[fileVal]
	if ok {
		n.Type = theType
	} else {
		n.Type = CNICTypes[0xffff]
	}

	return nil
}
