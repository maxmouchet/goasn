package goasn

import (
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

	db := ASNTree{
		treeV4: treeV4,
		treeV6: treeV6,
	}

	return &db, nil
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
