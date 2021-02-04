package block

type PartitionTableType string

type PartitionTable struct {
	Type       PartitionTableType
	Partitions []Partition
}
