package goasn

import (
	"errors"
	"net"

	"github.com/kentik/patricia"
	"github.com/kentik/patricia/uint32_tree"
)

type DB struct {
	treeV4 *uint32_tree.TreeV4
	treeV6 *uint32_tree.TreeV6
}

func NewDB(path string) (*DB, error) {
	treeV4 := uint32_tree.NewTreeV4()
	treeV6 := uint32_tree.NewTreeV6()

	entries, err := ReadPyASNDatabase(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		addrV4, addrV6, err := patricia.ParseFromIPAddr(&entry.IPNet)
		if err != nil {
			return nil, err
		}
		if addrV4 != nil {
			_, _, err = treeV4.Add(*addrV4, entry.ASN, nil)
		}
		if addrV6 != nil {
			_, _, err = treeV6.Add(*addrV6, entry.ASN, nil)
		}
		if err != nil {
			return nil, err
		}
	}

	db := DB{
		treeV4: treeV4,
		treeV6: treeV6,
	}

	return &db, nil
}

func (d *DB) LookupIP(ip net.IP) (uint32, error) {
	if addrV4 := ip.To4(); addrV4 != nil {
		addr := patricia.NewIPv4AddressFromBytes(addrV4, uint(32))
		_, asn, err := d.treeV4.FindDeepestTag(addr)
		return asn, err
	}
	if addrV6 := ip.To16(); addrV6 != nil {
		addr := patricia.NewIPv6Address(addrV6, 128)
		_, asn, err := d.treeV6.FindDeepestTag(addr)
		return asn, err
	}
	return 0, errors.New("Invalid IP address")
}

func (d *DB) LookupStr(s string) (uint32, error) {
	return d.LookupIP(net.ParseIP(s))
}
