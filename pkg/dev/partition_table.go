package dev

import (
	"path"
	"strings"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/partition/gpt"
	"github.com/diskfs/go-diskfs/partition/mbr"
	"github.com/pkg/errors"
)

const (
	CDevBasePath = "/dev"

	CGPTPartitionTableType = "GPT"
	CMBRPartitionTableType = "MBR"
)

type PartitionTableType string

type PartitionTable struct {
	Type       PartitionTableType
	Partitions []Partition
}

func NewPartitionTable(devName string) (*PartitionTable, error) {
	devPath := path.Join(CDevBasePath, devName)
	disk, err := diskfs.OpenWithMode(devPath, diskfs.ReadOnly)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open device %s", devPath)
	}

	recTable, err := disk.GetPartitionTable()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get disk partition table")
	}

	tableType := strings.ToUpper(recTable.Type())

	var table *PartitionTable
	switch tableType {
	case CGPTPartitionTableType:
		realTable := recTable.(*gpt.Table)
		table = &PartitionTable{
			Type:       CGPTPartitionTableType,
			Partitions: NewPartitionsFromGPT(realTable.Partitions),
		}
	case CMBRPartitionTableType:
		realTable := recTable.(*mbr.Table)
		table = &PartitionTable{
			Type:       CMBRPartitionTableType,
			Partitions: NewPartitionsFromMBR(realTable.Partitions),
		}
	default:
		return nil, errors.Errorf("unsupported partition table type %s", tableType)
	}

	return table, nil
}