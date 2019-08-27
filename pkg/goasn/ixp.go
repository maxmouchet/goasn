package goasn

import (
	"errors"
	"net"

	"github.com/kentik/patricia"
	"github.com/kentik/patricia/string_tree"
	"github.com/maxmouchet/goasn/pkg/peeringdb"
)

type IXPTree struct {
	treeV4 *string_tree.TreeV4
	treeV6 *string_tree.TreeV6
}

// PrefixIXP represents the IXP associated to a prefix.
type PrefixIXP struct {
	Prefix net.IPNet `json:"prefix"`
	IXP    string    `json:"ixp"`
}

func NewIXPTree(db peeringdb.DB) (*IXPTree, error) {
	treeV4 := string_tree.NewTreeV4()
	treeV6 := string_tree.NewTreeV6()

	ixsByID := make(map[int]peeringdb.IX)
	for _, ix := range db.IXs {
		ixsByID[ix.ID] = ix
	}

	ixLansByID := make(map[int]peeringdb.LAN)
	for _, ixlan := range db.LANs {
		ixLansByID[ixlan.ID] = ixlan
	}

	ixpfxByIX := make(map[int][]net.IPNet)
	for _, ixpfx := range db.Prefixes {
		lan := ixLansByID[ixpfx.IXLanID]
		_, ip, err := net.ParseCIDR(ixpfx.Prefix)
		if err != nil {
			return nil, err
		}
		ixpfxByIX[lan.IXID] = append(ixpfxByIX[lan.IXID], *ip)
	}

	for ixID, pfxs := range ixpfxByIX {
		ixp := ixsByID[ixID]
		for _, pfx := range pfxs {
			addrV4, addrV6, err := patricia.ParseFromIPAddr(&pfx)
			if err != nil {
				return nil, err
			}
			if addrV4 != nil {
				_, _, err = treeV4.Add(*addrV4, ixp.Name, nil)
			}
			if addrV6 != nil {
				_, _, err = treeV6.Add(*addrV6, ixp.Name, nil)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	tree := IXPTree{
		treeV4: treeV4,
		treeV6: treeV6,
	}

	return &tree, nil
}

// TODO: New IXPTreeFromFile
// Requires marshaling, ...

func (t IXPTree) LookupIP(ip net.IP) (string, error) {
	if addrV4 := ip.To4(); addrV4 != nil {
		addr := patricia.NewIPv4AddressFromBytes(addrV4, uint(32))
		_, ixp, err := t.treeV4.FindDeepestTag(addr)
		return ixp, err
	}
	if addrV6 := ip.To16(); addrV6 != nil {
		addr := patricia.NewIPv6Address(addrV6, 128)
		_, ixp, err := t.treeV6.FindDeepestTag(addr)
		return ixp, err
	}
	return "", errors.New("Invalid IP address")
}

func (t IXPTree) LookupIPMultiple(ips []net.IP) ([]string, error) {
	ixps := make([]string, len(ips))
	for i, ip := range ips {
		ixp, err := t.LookupIP(ip)
		if err != nil {
			return nil, err
		}
		ixps[i] = ixp
	}
	return ixps, nil
}

func (t IXPTree) LookupStr(str string) (string, error) {
	return t.LookupIP(net.ParseIP(str))
}

func (t IXPTree) LookupStrMultiple(strs []string) ([]string, error) {
	ixps := make([]string, len(strs))
	for i, str := range strs {
		ixp, err := t.LookupStr(str)
		if err != nil {
			return nil, err
		}
		ixps[i] = ixp
	}
	return ixps, nil
}
