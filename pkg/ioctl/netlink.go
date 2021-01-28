package ioctl

import (
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

type NetlinkSvc struct{}

func NewNetlinkSvc() *NetlinkSvc {
	return &NetlinkSvc{}
}

func (n *NetlinkSvc) GetIPv6NeighbourData() ([]IPv6Neighbour, error) {
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
			// TODO: continue instead of return
			return nil, errors.Wrapf(err, "unable to get neighbours for %s", iName)
		}

		for _, n := range nl {
			neighbour := NewIPv6Neighbour(iIdx, iName, &n)
			neighbours = append(neighbours, *neighbour)
		}
	}

	return neighbours, nil
}
