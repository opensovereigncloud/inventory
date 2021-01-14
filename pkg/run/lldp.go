package run

import (
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"net"
	"time"

	"github.com/mdlayher/lldp"
	"github.com/pkg/errors"
)

type LLDPFrameInfo struct {
	InterfaceID         string
	ChassisID           string
	SystemName          string
	SystemDescription   string
	Capabilities        []LLDPCapability
	EnabledCapabilities []LLDPCapability
	PortID              string
	PortDescription     string
	ManagementAddresses []string
	TTL                 time.Duration
}

func NewLLDPFrameInfo(interfaceID string, thePath string) (*LLDPFrameInfo, error) {
	contents, err := ioutil.ReadFile(thePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file %s", thePath)
	}

	// 1-8 bytes - ?
	// 9-22 bytes - ethernet frame part
	// 23-rest - LLDP frame part
	frame := lldp.Frame{}
	err = frame.UnmarshalBinary(contents[22:])
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal LLDP frame")
	}

	frameInfo := &LLDPFrameInfo{
		InterfaceID: interfaceID,
		TTL:         frame.TTL,
	}

	err = frameInfo.setChassisID(frame.ChassisID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set chassis ID")
	}

	err = frameInfo.setPortID(frame.PortID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal port ID")
	}

	for _, tlv := range frame.Optional {
		err = frameInfo.setOptional(tlv)
		if err != nil {
			return nil, errors.Wrap(err, "unable to process optional TLV")
		}
	}

	return frameInfo, nil
}

func (f *LLDPFrameInfo) setChassisID(chassisID *lldp.ChassisID) error {
	switch chassisID.Subtype {
	case lldp.ChassisIDSubtypeChassisComponenent:
		fallthrough
	case lldp.ChassisIDSubtypeInterfaceAlias:
		fallthrough
	case lldp.ChassisIDSubtypePortComponent:
		fallthrough
	case lldp.ChassisIDSubtypeInterfaceName:
		fallthrough
	case lldp.ChassisIDSubtypeLocallyAssigned:
		f.ChassisID = string(chassisID.ID)
	case lldp.ChassisIDSubtypeMACAddress:
		chassisId, err := idBytesToMac(chassisID.ID)
		if err != nil {
			return errors.Wrap(err, "unable to convert chassis ID to MAC address")
		}
		f.ChassisID = chassisId
	case lldp.ChassisIDSubtypeNetworkAddress:
		chassisId, err := idBytesToNetworkAddress(chassisID.ID)
		if err != nil {
			return errors.Wrap(err, "unable to convert chassis ID to MAC address")
		}
		f.ChassisID = chassisId
	default:
		// fallback
		f.ChassisID = hex.EncodeToString(chassisID.ID)
	}

	return nil
}

func (f *LLDPFrameInfo) setPortID(portID *lldp.PortID) error {
	switch portID.Subtype {
	case lldp.PortIDSubtypeInterfaceAlias:
		fallthrough
	case lldp.PortIDSubtypePortComponent:
		fallthrough
	case lldp.PortIDSubtypeInterfaceName:
		fallthrough
	case lldp.PortIDSubtypeLocallyAssigned:
		f.PortID = string(portID.ID)
	case lldp.PortIDSubtypeMACAddress:
		portId, err := idBytesToMac(portID.ID)
		if err != nil {
			return errors.Wrap(err, "unable to convert port ID to MAC address")
		}
		f.PortID = portId
	case lldp.PortIDSubtypeNetworkAddress:
		portId, err := idBytesToNetworkAddress(portID.ID)
		if err != nil {
			return errors.Wrap(err, "unable to convert port ID to network address")
		}
		f.PortID = portId
	default:
		// fallback
		f.PortID = hex.EncodeToString(portID.ID)
	}

	return nil
}

func (f *LLDPFrameInfo) setOptional(tlv *lldp.TLV) error {
	switch tlv.Type {
	case lldp.TLVTypeSystemName:
		f.setSystemName(tlv.Value)
	case lldp.TLVTypeSystemDescription:
		f.setSystemDescription(tlv.Value)
	case lldp.TLVTypePortDescription:
		f.setPortDescription(tlv.Value)
	case lldp.TLVTypeSystemCapabilities:
		f.setSystemCapability(tlv.Value)
	case lldp.TLVTypeManagementAddress:
		if err := f.setManagementAddress(tlv.Value); err != nil {
			return errors.Wrapf(err, "unable to set management address")
		}
		// there is also lldp.TLVTypeOrganizationSpecific
		// not collecting it for now
	}

	return nil
}

func (f *LLDPFrameInfo) setSystemName(val []byte) {
	f.SystemName = string(val)
}

func (f *LLDPFrameInfo) setSystemDescription(val []byte) {
	f.SystemDescription = string(val)
}

func (f *LLDPFrameInfo) setPortDescription(val []byte) {
	f.PortDescription = string(val)
}

func (f *LLDPFrameInfo) setSystemCapability(val []byte) {
	capabilities := binary.BigEndian.Uint16(val[0:2])
	enabledCapabilities := binary.BigEndian.Uint16(val[2:4])

	for i, capability := range CCapabilities {
		idx := uint16(1 << i)
		capable := capabilities & idx
		if capable != 0 {
			f.addCapability(capability)
		}

		enabled := enabledCapabilities & idx
		if enabled != 0 {
			f.addEnabledCapability(capability)
		}
	}
}

func (f *LLDPFrameInfo) addCapability(capability LLDPCapability) {
	if f.Capabilities == nil {
		f.Capabilities = make([]LLDPCapability, 0)
	}

	f.Capabilities = append(f.Capabilities, capability)
}

func (f *LLDPFrameInfo) addEnabledCapability(capability LLDPCapability) {
	if f.EnabledCapabilities == nil {
		f.EnabledCapabilities = make([]LLDPCapability, 0)
	}

	f.EnabledCapabilities = append(f.EnabledCapabilities, capability)
}

func (f *LLDPFrameInfo) setManagementAddress(val []byte) error {
	var addressBytes []byte

	valLen := len(val)
	if valLen < 6 {
		return errors.Errorf("expected to have at least 6 bytes, but got %d", valLen)
	}

	switch val[1] {
	case 1:
		addressBytes = val[2:6]
	case 2:
		addressBytes = val[2:18]
	default:
		return errors.Errorf("unhandled address type %d", val[1])
	}

	ip := net.IP(addressBytes)

	f.addManagementAddress(ip.String())

	return nil
}

func (f *LLDPFrameInfo) addManagementAddress(address string) {
	if f.ManagementAddresses == nil {
		f.ManagementAddresses = make([]string, 0)
	}

	f.ManagementAddresses = append(f.ManagementAddresses, address)
}
