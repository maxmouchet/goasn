package goasn

import (
	"net"

	"github.com/kentik/patricia/string_tree"
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
