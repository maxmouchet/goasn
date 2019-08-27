package goasn

import (
	"errors"
	"io/ioutil"
	"net"

	"github.com/kentik/patricia"
	"github.com/kentik/patricia/uint32_tree"
)

type ASNTree struct {
	treeV4 *uint32_tree.TreeV4
	treeV6 *uint32_tree.TreeV6
}

// PrefixOrigin represents the origin AS(es) associated to a prefix.
// There can be multiple ASes if the prefix is multi-homed.
type PrefixOrigin struct {
	Prefix net.IPNet `json:"prefix"`
	Origin []uint32  `json:"origin"`
}

func NewASNTree(prefixes []PrefixOrigin) (*ASNTree, error) {
	treeV4 := uint32_tree.NewTreeV4()
	treeV6 := uint32_tree.NewTreeV6()

	for _, prefix := range prefixes {
		addrV4, addrV6, err := patricia.ParseFromIPAddr(&prefix.Prefix)
		if err != nil {
			return nil, err
		}
		if addrV4 != nil {
			for _, asn := range prefix.Origin {
				_, _, err = treeV4.Add(*addrV4, asn, nil)
			}
		}
		if addrV6 != nil {
			for _, asn := range prefix.Origin {
				_, _, err = treeV6.Add(*addrV6, asn, nil)
			}
		}
		if err != nil {
			return nil, err
		}
	}

	tree := ASNTree{
		treeV4: treeV4,
		treeV6: treeV6,
	}

	return &tree, nil
}

// TODO: Detect file type (txt, json, ...)
func NewASNTreeFromFile(path string) (*ASNTree, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var db ASNDatabase
	err = db.UnmarshalText(b)
	if err != nil {
		return nil, err
	}

	return NewASNTree(db.Entries)
}

func NewPrefixOrigin(e RIBEntry) PrefixOrigin {
	var origin []uint32
	segment := e.Path[len(e.Path)-1]

	if segment.Type == ASPathSegmentTypeSequence {
		origin = []uint32{segment.ASNs[len(segment.ASNs)-1]}
	} else {
		origin = segment.ASNs
	}

	return PrefixOrigin{
		Prefix: e.Network,
		Origin: origin,
	}
}

func (t ASNTree) LookupIP(ip net.IP) ([]uint32, error) {
	if addrV4 := ip.To4(); addrV4 != nil {
		addr := patricia.NewIPv4AddressFromBytes(addrV4, uint(32))
		_, asn, err := t.treeV4.FindDeepestTags(addr)
		return asn, err
	}
	if addrV6 := ip.To16(); addrV6 != nil {
		addr := patricia.NewIPv6Address(addrV6, 128)
		_, asn, err := t.treeV6.FindDeepestTags(addr)
		return asn, err
	}
	return nil, errors.New("Invalid IP address")
}

func (t ASNTree) LookupIPMultiple(ips []net.IP) ([][]uint32, error) {
	asns := make([][]uint32, len(ips))
	for i, ip := range ips {
		asn, err := t.LookupIP(ip)
		if err != nil {
			return nil, err
		}
		asns[i] = asn
	}
	return asns, nil
}

func (t ASNTree) LookupStr(str string) ([]uint32, error) {
	return t.LookupIP(net.ParseIP(str))
}

func (t ASNTree) LookupStrMultiple(strs []string) ([][]uint32, error) {
	asns := make([][]uint32, len(strs))
	for i, str := range strs {
		asn, err := t.LookupStr(str)
		if err != nil {
			return nil, err
		}
		asns[i] = asn
	}
	return asns, nil
}
