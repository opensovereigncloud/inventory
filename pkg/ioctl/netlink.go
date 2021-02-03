package ioctl

import (
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"

	"github.com/onmetal/inventory/pkg/printer"
)

type NetlinkSvc struct {
	printer  *printer.Svc
	rootPath string
}

func NewNetlinkSvc(printer *printer.Svc, basePath string) *NetlinkSvc {
	return &NetlinkSvc{
		printer:  printer,
		rootPath: basePath,
	}
}

func (s *NetlinkSvc) GetIPv6NeighbourData() ([]IPv6Neighbour, error) {
	chroot, err := NewChroot(s.rootPath)
	defer func() {
		if err := chroot.Close(); err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to exit chroot"))
		}
	}()

	ll, err := netlink.LinkList()
	if err != nil {
		return nil, errors.Wrap(err, "unable to obtain device list")
	}

	neighbours := make([]IPv6Neighbour, 0)
	for _, l := range ll {
		iIdx := l.Attrs().Index
		iName := l.Attrs().Name
		nl, err := netlink.NeighList(iIdx, netlink.FAMILY_V6)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to get neighbours for %s", iName))
			continue
		}

		for _, n := range nl {
			neighbour := NewIPv6Neighbour(iIdx, iName, &n)
			neighbours = append(neighbours, *neighbour)
		}
	}

	return neighbours, nil
}
