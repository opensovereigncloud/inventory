package ipmi

import (
	"fmt"
	"net"
	"os"

	"github.com/pkg/errors"
	"github.com/u-root/u-root/pkg/ipmi"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CIPMIIOCtlChannel = 1

	CIPMIIOCtlSetInProgressFlag   = 0
	CIPMIIOCtlIPAddressFlag       = 3
	CIPMIIOCtlIPAddressSourceFlag = 4
	CIPMIIOCtlMACAddressFlag      = 5

	CIPMISetInProgressSetCompleteStatus   = "Set Complete"
	CIPMISetInProgressSetInProgressStatus = "Set In Progress"
	CIPMISetInProgressCommitWriteStatus   = "Commit Write"
	CIPMISetInProgressReservedStatus      = "Reserved"
	CIPMISetInProgressUnknownStatus       = "Unknown"

	CIPMIIPAddressSourceUnspecified = "Unspecified"
	CIPMIIPAddressSourceStatic      = "Static Address"
	CIPMIIPAddressSourceDHCP        = "DHCP Address"
	CIPMIIPAddressSourceBIOS        = "BIOS Assigned Address"
	CIPMIIPAddressSourceOther       = "Other"

	CIPMIAdditionalSensorDeviceSupport        = "Sensor Device"
	CIPMIAdditionalSDRRepositoryDeviceSupport = "SDR Repository Device"
	CIPMIAdditionalSELDeviceSupport           = "SEL Device"
	CIPMIAdditionalFRUInventoryDeviceSupport  = "FRU Inventory Device"
	CIPMIAdditionalIPMBEventReceiverSupport   = "IPMB Event Receiver"
	CIPMIAdditionalIPMBEventGeneratorSupport  = "IPMB Event Generator"
	CIPMIAdditionalBridgeSupport              = "Bridge"
	CIPMIAdditionalChassisDeviceSupport       = "Chassis Device"
)

type IPMISetInProgressStatus string

var CIPMISetInProgressStatuses = []IPMISetInProgressStatus{
	CIPMISetInProgressSetCompleteStatus,
	CIPMISetInProgressSetInProgressStatus,
	CIPMISetInProgressCommitWriteStatus,
	CIPMISetInProgressReservedStatus,
}

type IPMIIPAddressSource string

var CIPMIIPAddressSources = []IPMIIPAddressSource{
	CIPMIIPAddressSourceUnspecified,
	CIPMIIPAddressSourceStatic,
	CIPMIIPAddressSourceDHCP,
	CIPMIIPAddressSourceBIOS,
}

type IPMIAdditionalDeviceSupport string

var CIPMIAdditionalDeviceSupportList = []IPMIAdditionalDeviceSupport{
	CIPMIAdditionalSensorDeviceSupport,
	CIPMIAdditionalSDRRepositoryDeviceSupport,
	CIPMIAdditionalSELDeviceSupport,
	CIPMIAdditionalFRUInventoryDeviceSupport,
	CIPMIAdditionalIPMBEventReceiverSupport,
	CIPMIAdditionalIPMBEventGeneratorSupport,
	CIPMIAdditionalBridgeSupport,
	CIPMIAdditionalChassisDeviceSupport,
}

type IPMIDeviceInfo struct {
	ID                      uint8
	Revision                uint8
	FirmwareRevision        string
	IPMIVersion             string
	ManufacturerID          string
	ProductID               string
	DeviceAvailable         bool
	ProvidesDeviceSDRs      bool
	AdditionalDeviceSupport []IPMIAdditionalDeviceSupport
	AuxFirmwareRevInfo      []string

	SetInProgress   IPMISetInProgressStatus
	IPAddressSource IPMIIPAddressSource
	IPAddress       string
	MACAddress      string
}

type IPMIDeviceInfoSvc struct {
	printer *printer.Svc
}

func NewIPMIDeviceInfoSvc(printer *printer.Svc) *IPMIDeviceInfoSvc {
	return &IPMIDeviceInfoSvc{
		printer: printer,
	}
}

func (s *IPMIDeviceInfoSvc) GetIPMIDeviceInfo(thePath string) (*IPMIDeviceInfo, error) {
	f, err := os.OpenFile(thePath, os.O_RDWR, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open IPMI device file %s", thePath)
	}

	conn := &ipmi.IPMI{
		File: f,
	}
	defer func() {
		if err := conn.Close(); err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to close file %s", thePath))
		}
	}()

	info := &IPMIDeviceInfo{}

	defs := []func(*ipmi.IPMI) error{
		info.defDevice,
		info.defSetInProgress,
		info.defIPAddressSource,
		info.defIPAddress,
		info.defMACAddress,
	}

	for _, def := range defs {
		err := def(conn)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to set IPMI device info property"))
		}
	}

	return info, nil
}

func (i *IPMIDeviceInfo) defDevice(conn *ipmi.IPMI) error {
	deviceInfo, err := conn.GetDeviceID()
	if err != nil {
		return err
	}

	i.ID = deviceInfo.DeviceID
	i.Revision = deviceInfo.DeviceRevision & 0x0f
	i.FirmwareRevision = fmt.Sprintf("%d.%02x", deviceInfo.FwRev1&0x3f, deviceInfo.FwRev2)
	i.IPMIVersion = fmt.Sprintf("%x.%x", deviceInfo.IpmiVersion&0x0f, (deviceInfo.IpmiVersion&0x0f)>>4)

	if deviceInfo.AdtlDeviceSupport != 0 {
		i.AdditionalDeviceSupport = make([]IPMIAdditionalDeviceSupport, 0)
	}

	manufacturerID := uint32(deviceInfo.ManufacturerID[2]) << 16
	manufacturerID |= uint32(deviceInfo.ManufacturerID[1]) << 8
	manufacturerID |= uint32(deviceInfo.ManufacturerID[0])
	i.ManufacturerID = fmt.Sprintf("%d (0x%04X)", manufacturerID, manufacturerID)

	productID := uint16(deviceInfo.ProductID[1]) << 8
	productID |= uint16(deviceInfo.ProductID[1])
	i.ProductID = fmt.Sprintf("%d (0x%04X)", productID, productID)

	i.DeviceAvailable = (^deviceInfo.FwRev1 & 0x80) != 0
	i.ProvidesDeviceSDRs = deviceInfo.DeviceRevision&0x80 != 0

	for shift, sup := range CIPMIAdditionalDeviceSupportList {
		idx := byte(1 << shift)
		val := deviceInfo.AdtlDeviceSupport & idx

		if val != 0 {
			i.AdditionalDeviceSupport = append(i.AdditionalDeviceSupport, sup)
		}
	}

	i.AuxFirmwareRevInfo = make([]string, len(deviceInfo.AuxFwRev))

	for idx, val := range deviceInfo.AuxFwRev {
		i.AuxFirmwareRevInfo[idx] = fmt.Sprintf("0x%02x", val)
	}

	return nil
}

func (i *IPMIDeviceInfo) defSetInProgress(conn *ipmi.IPMI) error {
	bytes, err := conn.GetLanConfig(CIPMIIOCtlChannel, CIPMIIOCtlSetInProgressFlag)
	if err != nil {
		return errors.Wrap(err, "unable to get set in progress")
	}

	if len(bytes) < 3 {
		return errors.Wrap(err, "unable to get set in progress")
	}

	valIdx := int(bytes[2])

	if valIdx >= len(CIPMISetInProgressStatuses) {
		i.SetInProgress = CIPMISetInProgressUnknownStatus
	}

	i.SetInProgress = CIPMISetInProgressStatuses[valIdx]

	return nil
}

func (i *IPMIDeviceInfo) defIPAddressSource(conn *ipmi.IPMI) error {
	bytes, err := conn.GetLanConfig(CIPMIIOCtlChannel, CIPMIIOCtlIPAddressSourceFlag)
	if err != nil {
		return errors.Wrap(err, "unable to get IP address source")
	}

	if len(bytes) < 3 {
		return errors.New("unable to obtain IP address source")
	}

	valIdx := int(bytes[2])

	if valIdx >= len(CIPMIIPAddressSources) {
		i.IPAddressSource = CIPMIIPAddressSourceOther
	}

	i.IPAddressSource = CIPMIIPAddressSources[valIdx]

	return nil
}

func (i *IPMIDeviceInfo) defIPAddress(conn *ipmi.IPMI) error {
	bytes, err := conn.GetLanConfig(CIPMIIOCtlChannel, CIPMIIOCtlIPAddressFlag)
	if err != nil {
		return errors.Wrap(err, "unable to obtain IP address")
	}

	if len(bytes) < 6 {
		return errors.New("unable to obtain IP address")
	}

	ipBytes := bytes[2:6]
	ip := net.IP(ipBytes)
	i.IPAddress = ip.String()

	return nil
}

func (i *IPMIDeviceInfo) defMACAddress(conn *ipmi.IPMI) error {
	bytes, err := conn.GetLanConfig(CIPMIIOCtlChannel, CIPMIIOCtlMACAddressFlag)
	if err != nil {
		return errors.Wrap(err, "unable to obtain MAC address")
	}

	if len(bytes) < 8 {
		return errors.New("unable to obtain MAC address")
	}

	macBytes := bytes[2:8]
	mac := net.HardwareAddr(macBytes)
	i.MACAddress = mac.String()

	return nil
}
