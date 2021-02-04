package pci

type Bus struct {
	ID      string
	Devices []Device
}

func NewBus(id string, devices []Device) *Bus {
	return &Bus{
		ID:      id,
		Devices: devices,
	}
}
