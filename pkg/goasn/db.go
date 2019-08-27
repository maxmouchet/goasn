package goasn

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

// pyasn-style database

type ASNDatabase []PrefixOrigin

type IXPDatabase []PrefixIXP

// TODO: Deduplicate code
func (db ASNDatabase) MarshalText(singleAS bool) ([]byte, error) {
	w := new(bytes.Buffer)

	// TODO: Diff with pyasn
	// TODO: Original source, Prefixes-v4/v6
	fmt.Fprintf(w, "; IP-ASN32-DAT file\n")
	fmt.Fprintf(w, "; Original source:\n")
	fmt.Fprintf(w, "; Converted on:\t%s\n", time.Now().Format("Mon Jan 2 15:04:05 2006"))
	fmt.Fprintf(w, "; Prefixes-v4:\t%d\n")
	fmt.Fprintf(w, "; Prefixes-v6:\t%d\n")
	fmt.Fprintf(w, ";\n")

	lastNet := ""

	_, defaultV4, _ := net.ParseCIDR("0.0.0.0/0")
	_, defaultV6, _ := net.ParseCIDR("::/0")

	// WARN if same prefix with differents ASes

	for _, entry := range db {
		if entry.Prefix.String() == lastNet {
			continue
		}

		// TODO: Optimize
		if (entry.Prefix.String() == defaultV4.String()) || (entry.Prefix.String() == defaultV6.String()) {
			continue
		}

		lastNet = entry.Prefix.String()
		asns := entry.Origin
		if singleAS {
			asns = asns[0:1]
		}
		fmt.Fprintf(
			w,
			"%s\t%s\n",
			entry.Prefix.String(),
			formatSlice(asns),
		)
	}

	return w.Bytes(), nil
}

// func (db *ASNDatabase) UnmarshalText(data []byte) error {

// }

func formatSlice(s []uint32) string {
	if len(s) == 0 {
		return ""
	}
	str := fmt.Sprintf("%d", s[0])
	for _, e := range s[1:] {
		str += fmt.Sprintf(",%d", e)
	}
	return str
}
