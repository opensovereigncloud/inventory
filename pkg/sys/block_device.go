package sys

import (
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/file"
)

const (
	CQueueRotationalPath        = "/queue/rotational"
	CQueuePhysicalBlockSizePath = "/queue/physical_block_size"
	CDeviceVendorPath           = "/device/vendor"
	CDeviceModelPath            = "/device/model"
	CDeviceSerialPath           = "/device/serial"
	CDeviceNumaNodePath         = "/device/numa_node"
	CDeviceFirmwareRevPath      = "/device/firmware_rev"
	CDeviceState                = "/device/state"
	CWWIDPath                   = "/wwid"
	CRemovablePath              = "/removable"
	CSizePath                   = "/size"
	CReadOnlyPath               = "/ro"

	CDefaultSectorSize = 512

	CFloppyDiskPattern = "fd(\\d+)"
	CCDROMDiskPattern  = "(sr|scd)(\\d+)"
	CSCSIDiskPattern   = "sd(\\w+)"
	CIDEDiskPattern    = "hd(\\w+)"
	CNVMeDiskPrefix    = "nvme(\\d+)(n\\d+)?(p\\d+)?"
	CUSBDiskPattern    = "ub(\\w+)"
	CMMCDiskPattern    = "mmcblk(\\d+)(p\\d+)?"
	CVirtIODiskPattern = "vd(\\w+)"
	CXenDiskPattern    = "xvd(\\w+)"

	CFloppyDiskName = "Floppy"
	CCDROMDiskName  = "CD-ROM"
	CSCSIDiskName   = "SCSI"
	CIDEDiskName    = "IDE"
	CNVMeDiskName   = "NVMe"
	CUSBDiskName    = "USB"
	CMMCDiskName    = "MMC"
	CVirtIODiskName = "VirtIO"
	CXenDiskName    = "Xen"
)

var CFloppyDiskRegexp = regexp.MustCompile(CFloppyDiskPattern)
var CCDROMDiskRegexp = regexp.MustCompile(CCDROMDiskPattern)
var CSCSIDiskRegexp = regexp.MustCompile(CSCSIDiskPattern)
var CIDEDiskRegexp = regexp.MustCompile(CIDEDiskPattern)
var CNVMeDiskRegexp = regexp.MustCompile(CNVMeDiskPrefix)
var CUSBDiskRegexp = regexp.MustCompile(CUSBDiskPattern)
var CMMCDiskRegexp = regexp.MustCompile(CMMCDiskPattern)
var CVirtIODiskRegexp = regexp.MustCompile(CVirtIODiskPattern)
var CXenDiskRegexp = regexp.MustCompile(CXenDiskPattern)

var CDiskMap = map[string]regexp.Regexp{
	CFloppyDiskName: *CFloppyDiskRegexp,
	CCDROMDiskName:  *CCDROMDiskRegexp,
	CSCSIDiskName:   *CSCSIDiskRegexp,
	CIDEDiskName:    *CIDEDiskRegexp,
	CNVMeDiskName:   *CNVMeDiskRegexp,
	CUSBDiskName:    *CUSBDiskRegexp,
	CMMCDiskName:    *CMMCDiskRegexp,
	CVirtIODiskName: *CVirtIODiskRegexp,
	CXenDiskName:    *CXenDiskRegexp,
}

type BlockDevice struct {
	Path             string
	Name             string
	Type             string
	Rotational       bool
	Removable        bool
	ReadOnly         bool
	Vendor           string
	Model            string
	Serial           string
	WWID             string
	FirmwareRevision string
	State            string
	Size             uint64
	BlockSize        uint64
	NUMANodeID       uint64
	Stat             *BlockDeviceStat
}

func NewBlockDevice(thePath string, name string) (*BlockDevice, error) {
	device := &BlockDevice{
		Path: thePath,
		Name: name,
	}

	defs := []func() error{
		device.defType,
		device.defRotational,
		device.defRemovable,
		device.defReadOnly,
		device.defVendor,
		device.defModel,
		device.defSerial,
		device.defWWID,
		device.defFirmwareRev,
		device.defState,
		device.defSize,
		device.defBlockSize,
		device.defNumaNodeID,
		device.defStat,
	}

	errs := make([]error, 0)

	for _, def := range defs {
		err := def()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return device, nil
}

func (bd *BlockDevice) defType() error {
	for k, v := range CDiskMap {
		if v.MatchString(bd.Name) {
			bd.Type = k
			break
		}
	}

	if bd.Type == "" {
		return errors.Errorf("unhandled block device type %s found", bd.Name)
	}

	return nil
}

func (bd *BlockDevice) defRotational() error {
	rotationalPath := path.Join(bd.Path, CQueueRotationalPath)
	rotational, err := file.ToBool(rotationalPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", rotationalPath)
	}

	bd.Rotational = rotational

	return nil
}

func (bd *BlockDevice) defVendor() error {
	vendorPath := path.Join(bd.Path, CDeviceVendorPath)
	vendor, err := file.ToString(vendorPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", vendorPath)
	}

	bd.Vendor = vendor

	return nil
}

func (bd *BlockDevice) defModel() error {
	modelPath := path.Join(bd.Path, CDeviceModelPath)
	model, err := file.ToString(modelPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", modelPath)
	}

	bd.Model = model

	return nil
}

func (bd *BlockDevice) defWWID() error {
	wwidPath := path.Join(bd.Path, CWWIDPath)
	wwid, err := file.ToString(wwidPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", wwidPath)
	}

	bd.WWID = wwid

	return nil
}

func (bd *BlockDevice) defRemovable() error {
	removablePath := path.Join(bd.Path, CRemovablePath)
	removable, err := file.ToBool(removablePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", removablePath)
	}

	bd.Removable = removable

	return nil
}

func (bd *BlockDevice) defSerial() error {
	serialPath := path.Join(bd.Path, CDeviceSerialPath)
	serial, err := file.ToString(serialPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", serialPath)
	}

	bd.Serial = serial

	return nil
}

func (bd *BlockDevice) defSize() error {
	sizePath := path.Join(bd.Path, CSizePath)
	size, err := file.ToUint64(sizePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", sizePath)
	}

	bd.Size = size * CDefaultSectorSize

	return nil
}

func (bd *BlockDevice) defBlockSize() error {
	blockSizePath := path.Join(bd.Path, CQueuePhysicalBlockSizePath)
	blockSize, err := file.ToUint64(blockSizePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", blockSizePath)
	}

	bd.BlockSize = blockSize

	return nil
}

func (bd *BlockDevice) defNumaNodeID() error {
	numaPath := path.Join(bd.Path, CDeviceNumaNodePath)
	numa, err := file.ToUint64(numaPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", numaPath)
	}

	bd.NUMANodeID = numa

	return nil
}

func (bd *BlockDevice) defReadOnly() error {
	roPath := path.Join(bd.Path, CReadOnlyPath)
	ro, err := file.ToBool(roPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", roPath)
	}

	bd.ReadOnly = ro

	return nil
}

func (bd *BlockDevice) defFirmwareRev() error {
	fwPath := path.Join(bd.Path, CDeviceFirmwareRevPath)
	fw, err := file.ToString(fwPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", fwPath)
	}

	bd.FirmwareRevision = fw

	return nil
}

func (bd *BlockDevice) defState() error {
	statePath := path.Join(bd.Path, CDeviceState)
	state, err := file.ToString(statePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", statePath)
	}

	bd.State = state

	return nil
}

func (bd *BlockDevice) defStat() error {
	stat, err := NewBlockDeviceStat(bd.Path)
	if err != nil {
		return errors.Wrap(err, "unable to collect stats")
	}

	bd.Stat = stat

	return nil
}
