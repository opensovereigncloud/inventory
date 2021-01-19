package sys

import (
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/file"
	"github.com/onmetal/inventory/pkg/pci"
)

const (
	CPCIDeviceVendorPath          = "/vendor"
	CPCIDeviceTypePath            = "/device"
	CPCIDeviceSubsystemDevicePath = "/subsystem_device"
	CPCIDeviceSubsystemVendorPath = "/subsystem_vendor"

	CPCIDeviceClassPath = "/class"

	CPCIDeviceClassPattern = "0x([[:xdigit:]]{2})([[:xdigit:]]{2})([[:xdigit:]]{2})"
)

var CPCIDeviceClassRegexp = regexp.MustCompile(CPCIDeviceClassPattern)

type PCIDeviceType struct {
	ID   string
	Name string
}

type PCIDeviceVendor struct {
	ID   string
	Name string
}

type PCIDeviceSubtype struct {
	ID   string
	Name string
}

type PCIDeviceClass struct {
	ID   string
	Name string
}

type PCIDeviceSubclass struct {
	ID   string
	Name string
}

type PCIDeviceProgrammingInterface struct {
	ID   string
	Name string
}

type PCIDevice struct {
	Address              string
	Vendor               *PCIDeviceVendor
	Type                 *PCIDeviceType
	Subvendor            *PCIDeviceVendor
	Subtype              *PCIDeviceSubtype
	Class                *PCIDeviceClass
	Subclass             *PCIDeviceSubclass
	ProgrammingInterface *PCIDeviceProgrammingInterface
}

func NewPCIDevice(thePath string, name string, ids *pci.IDs) (*PCIDevice, error) {
	device := &PCIDevice{
		Address: name,
	}

	err := device.defVendor(thePath, ids)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set vendor branch")
	}

	err = device.defClass(thePath, ids)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set class branch")
	}

	return device, nil
}

func (d *PCIDevice) defVendor(thePath string, ids *pci.IDs) error {
	vendorPath := path.Join(thePath, CPCIDeviceVendorPath)
	vendorString, err := file.ToString(vendorPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get vendor string from %s", vendorPath)
	}

	vendorVal := vendorString[2:]
	vendor, ok := ids.Vendors[vendorVal]
	if !ok {
		// return errors.Errorf("unknown vendor id %s", vendorVal)
		return nil
	}

	d.Vendor = &PCIDeviceVendor{
		ID:   vendor.ID,
		Name: vendor.Name,
	}

	if len(vendor.Devices) > 0 {
		return d.defType(thePath, &vendor, ids)
	}

	return nil
}

func (d *PCIDevice) defType(thePath string, vendor *pci.Vendor, ids *pci.IDs) error {
	typePath := path.Join(thePath, CPCIDeviceTypePath)
	typeString, err := file.ToString(typePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get device/type string from %s", typePath)
	}

	typeVal := typeString[2:]
	theType, ok := vendor.Devices[typeVal]
	if !ok {
		// return errors.Errorf("unknown device/type id %s", typeVal)
		return nil
	}

	d.Type = &PCIDeviceType{
		ID:   theType.ID,
		Name: theType.Name,
	}

	if len(theType.Subsystems) > 0 {
		return d.defSubsystem(thePath, &theType, ids)
	}

	return nil
}

func (d *PCIDevice) defSubsystem(thePath string, theType *pci.Device, ids *pci.IDs) error {
	subDevicePath := path.Join(thePath, CPCIDeviceSubsystemDevicePath)
	subDeviceString, err := file.ToString(subDevicePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get subsystem device/type string from %s", subDevicePath)
	}
	subDeviceVal := subDeviceString[2:]

	subVendorPath := path.Join(thePath, CPCIDeviceSubsystemVendorPath)
	subVendorString, err := file.ToString(subVendorPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get subsystem vendor string from %s", subVendorPath)
	}
	subVendorVal := subVendorString[2:]

	subsystem, ok := theType.Subsystems[subVendorVal+subDeviceVal]
	if !ok {
		// TODO think how to log undiscovered items and report them to PCI IDs maintainers
		// return errors.Errorf("unknown subsystem id %s %s", subVendorVal, subDeviceVal)
		return nil
	}

	d.Subtype = &PCIDeviceSubtype{
		ID:   subsystem.SubdeviceID,
		Name: subsystem.Name,
	}

	subvendor, ok := ids.Vendors[subVendorVal]
	if !ok {
		// return errors.Errorf("unknown vendor id %s", subVendorVal)
		return nil
	}

	d.Subvendor = &PCIDeviceVendor{
		ID:   subvendor.ID,
		Name: subvendor.Name,
	}

	return nil
}

func (d *PCIDevice) defClass(thePath string, ids *pci.IDs) error {
	classPath := path.Join(thePath, CPCIDeviceClassPath)
	classString, err := file.ToString(classPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get vendor string from %s", classPath)
	}

	groups := CPCIDeviceClassRegexp.FindStringSubmatch(classString)
	classVal := groups[1]
	subclassVal := groups[2]
	ifaceVal := groups[3]

	class, ok := ids.Classes[classVal]
	if ok {
		d.Class = &PCIDeviceClass{
			ID:   class.ID,
			Name: class.Name,
		}
	}

	if len(class.Subclasses) == 0 {
		return nil
	}

	subclass, ok := class.Subclasses[subclassVal]
	if ok {
		d.Subclass = &PCIDeviceSubclass{
			ID:   subclass.ID,
			Name: subclass.Name,
		}
	}

	if len(subclass.ProgrammingInterfaces) == 0 {
		return nil
	}

	iface, ok := subclass.ProgrammingInterfaces[ifaceVal]
	if ok {
		d.ProgrammingInterface = &PCIDeviceProgrammingInterface{
			ID:   iface.ID,
			Name: iface.Name,
		}
	}

	return nil
}
