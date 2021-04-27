package host

import (
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/utils"
	"github.com/pkg/errors"
	"os"
)

type Info struct {
	Type string
	Name string
}

type Svc struct {
	printer *printer.Svc
}

func NewSvc(printer *printer.Svc) *Svc {
	return &Svc{
		printer: printer,
	}
}

func (s *Svc) GetData() (*Info, error) {
	info := Info{}
	name, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hostname")
	}
	info.Name = name
	hostType, err := utils.GetHostType()
	if err != nil {
		return nil, errors.Wrap(err, "failed to determine host type")
	}
	info.Type = hostType
	return &info, nil
}
