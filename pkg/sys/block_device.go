package sys

import (
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/dev"
	"github.com/onmetal/inventory/pkg/file"
	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CQueueRotationalPath        = "/queue/rotational"
	CQueuePhysicalBlockSizePath = "/queue/physical_block_size"
	CQueueLogicalBlockSizePath  = "/queue/logical_block_size"
	CQueueHWSectorSizePath      = "/queue/hw_sector_size"
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
	path              string
	Name              string
	Type              string
	Rotational        bool
	Removable         bool
	ReadOnly          bool
	Vendor            string
	Model             string
	Serial            string
	WWID              string
	FirmwareRevision  string
	State             string
	PhysicalBlockSize uint64
	LogicalBlockSize  uint64
	HWSectorSize      uint64
	Size              uint64
	NUMANodeID        uint64
	PartitionTable    *dev.PartitionTable
	Stat              *BlockDeviceStat
}

type BlockDeviceSvc struct {
	printer            *printer.Svc
	partitionTableSvc  *dev.PartitionTableSvc
	blockDeviceStatSvc *BlockDeviceStatSvc
}

func NewBlockDeviceSvc(printer *printer.Svc, partTableSvc *dev.PartitionTableSvc, blockDeviceStatSvc *BlockDeviceStatSvc) *BlockDeviceSvc {
	return &BlockDeviceSvc{
		printer:            printer,
		partitionTableSvc:  partTableSvc,
		blockDeviceStatSvc: blockDeviceStatSvc,
	}
}

func (s *BlockDeviceSvc) GetBlockDevice(thePath string, name string) (*BlockDevice, error) {
	device := &BlockDevice{
		path: thePath,
		Name: name,
	}

	defs := []func(*BlockDevice) error{
		s.defType,
		s.defRotational,
		s.defRemovable,
		s.defReadOnly,
		s.defVendor,
		s.defModel,
		s.defSerial,
		s.defWWID,
		s.defFirmwareRev,
		s.defState,
		s.defHWSectorSize,
		s.defPhysicalBlockSize,
		s.defLogicalBlockSize,
		s.defSize,
		s.defNumaNodeID,
		s.defPartitionTable,
		s.defStat,
	}

	for _, def := range defs {
		err := def(device)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to set block device property"))
		}
	}

	return device, nil
}

func (s *BlockDeviceSvc) defType(bd *BlockDevice) error {
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

func (s *BlockDeviceSvc) defRotational(bd *BlockDevice) error {
	rotationalPath := path.Join(bd.path, CQueueRotationalPath)
	rotational, err := file.ToBool(rotationalPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", rotationalPath)
	}

	bd.Rotational = rotational

	return nil
}

func (s *BlockDeviceSvc) defVendor(bd *BlockDevice) error {
	vendorPath := path.Join(bd.path, CDeviceVendorPath)
	vendor, err := file.ToString(vendorPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", vendorPath)
	}

	bd.Vendor = vendor

	return nil
}

func (s *BlockDeviceSvc) defModel(bd *BlockDevice) error {
	modelPath := path.Join(bd.path, CDeviceModelPath)
	model, err := file.ToString(modelPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", modelPath)
	}

	bd.Model = model

	return nil
}

func (s *BlockDeviceSvc) defWWID(bd *BlockDevice) error {
	wwidPath := path.Join(bd.path, CWWIDPath)
	wwid, err := file.ToString(wwidPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", wwidPath)
	}

	bd.WWID = wwid

	return nil
}

func (s *BlockDeviceSvc) defRemovable(bd *BlockDevice) error {
	removablePath := path.Join(bd.path, CRemovablePath)
	removable, err := file.ToBool(removablePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", removablePath)
	}

	bd.Removable = removable

	return nil
}

func (s *BlockDeviceSvc) defSerial(bd *BlockDevice) error {
	serialPath := path.Join(bd.path, CDeviceSerialPath)
	serial, err := file.ToString(serialPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", serialPath)
	}

	bd.Serial = serial

	return nil
}

func (s *BlockDeviceSvc) defSize(bd *BlockDevice) error {
	sizePath := path.Join(bd.path, CSizePath)
	size, err := file.ToUint64(sizePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", sizePath)
	}

	bd.Size = size * bd.HWSectorSize

	return nil
}

func (s *BlockDeviceSvc) defPhysicalBlockSize(bd *BlockDevice) error {
	blockSizePath := path.Join(bd.path, CQueuePhysicalBlockSizePath)
	blockSize, err := file.ToUint64(blockSizePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", blockSizePath)
	}

	bd.PhysicalBlockSize = blockSize

	return nil
}

func (s *BlockDeviceSvc) defLogicalBlockSize(bd *BlockDevice) error {
	blockSizePath := path.Join(bd.path, CQueueLogicalBlockSizePath)
	blockSize, err := file.ToUint64(blockSizePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", blockSizePath)
	}

	bd.LogicalBlockSize = blockSize

	return nil
}

func (s *BlockDeviceSvc) defHWSectorSize(bd *BlockDevice) error {
	sectorSizePath := path.Join(bd.path, CQueueHWSectorSizePath)
	sectorSize, err := file.ToUint64(sectorSizePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", sectorSizePath)
	}

	bd.HWSectorSize = sectorSize

	return nil
}

func (s *BlockDeviceSvc) defNumaNodeID(bd *BlockDevice) error {
	numaPath := path.Join(bd.path, CDeviceNumaNodePath)
	numa, err := file.ToUint64(numaPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", numaPath)
	}

	bd.NUMANodeID = numa

	return nil
}

func (s *BlockDeviceSvc) defReadOnly(bd *BlockDevice) error {
	roPath := path.Join(bd.path, CReadOnlyPath)
	ro, err := file.ToBool(roPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", roPath)
	}

	bd.ReadOnly = ro

	return nil
}

func (s *BlockDeviceSvc) defFirmwareRev(bd *BlockDevice) error {
	fwPath := path.Join(bd.path, CDeviceFirmwareRevPath)
	fw, err := file.ToString(fwPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", fwPath)
	}

	bd.FirmwareRevision = fw

	return nil
}

func (s *BlockDeviceSvc) defState(bd *BlockDevice) error {
	statePath := path.Join(bd.path, CDeviceState)
	state, err := file.ToString(statePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", statePath)
	}

	bd.State = state

	return nil
}

func (s *BlockDeviceSvc) defPartitionTable(bd *BlockDevice) error {
	table, err := s.partitionTableSvc.GetPartitionTable(bd.Name)
	if err != nil {
		return errors.Wrap(err, "unable to get partition table")
	}

	bd.PartitionTable = table

	return nil
}

func (s *BlockDeviceSvc) defStat(bd *BlockDevice) error {
	stat, err := s.blockDeviceStatSvc.GetBlockDeviceStat(bd.path)
	if err != nil {
		return errors.Wrap(err, "unable to collect stats")
	}

	bd.Stat = stat

	return nil
}
