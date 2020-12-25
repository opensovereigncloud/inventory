package sys

type PCIBus struct {
	ID string
}

func NewPCIBus(thePath string, id string) (*PCIBus, error) {
	return &PCIBus{
		ID: id,
	}, nil
}
