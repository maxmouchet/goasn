package goasn

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"github.com/kaorimatz/go-mrt"
)

// MRT format

func ribEntryFromMRT(entry mrt.TableDumpV2RIBEntry, network net.IPNet) RIBEntry {
	ribEntry := RIBEntry{}
	ribEntry.Network = network
	ribEntry.Path = make([]uint32, 0)

	for _, attr := range entry.BGPAttributes {
		// http://www.networkers-online.com/blog/2012/05/bgp-attributes/
		switch attr.TypeCode {
		case 1:
			origin := attr.Value.(mrt.BGPPathAttributeOrigin)
			switch origin {
			case mrt.BGPPathAttributeOriginIGP:
				ribEntry.OriginCode = "i"
			case mrt.BGPPathAttributeOriginEGP:
				ribEntry.OriginCode = "e"
			case mrt.BGPPathAttributeOriginIncomplete:
				ribEntry.OriginCode = "?"
			}

		case 2:
			asPath := attr.Value.(mrt.BGPPathAttributeASPath)
			for _, asPathSegment := range asPath {
				// TODO: Differentiate set from sequence
				// asPathSegment.Type
				for _, as := range asPathSegment.Value {
					var asn uint32
					switch len(as) {
					case 2:
						asn = uint32(binary.BigEndian.Uint16(as))
					case 4:
						asn = binary.BigEndian.Uint32(as)
					}
					// TODO: Is this correct ?
					ribEntry.Path = append(ribEntry.Path, asn)
				}
			}

		case 5:
			locPrf := attr.Value.(uint32)
			ribEntry.LocPrf = locPrf
		}
	}

	return ribEntry
}

func ribEntriesFromMRT(r io.Reader) ([]RIBEntry, error) {
	mrtReader := mrt.NewReader(r)
	ribEntries := make([]RIBEntry, 0)
	i := 0

	for {
		record, err := mrtReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch (*record).Type() {
		case mrt.TYPE_TABLE_DUMP_V2:
			switch (*record).Subtype() {
			case
				mrt.TABLE_DUMP_V2_SUBTYPE_RIB_IPv4_UNICAST,
				mrt.TABLE_DUMP_V2_SUBTYPE_RIB_IPv6_UNICAST:
				rib := (*record).(*mrt.TableDumpV2RIB)
				for _, entry := range rib.RIBEntries {
					ribEntries = append(ribEntries, ribEntryFromMRT(*entry, *rib.Prefix))
				}
			}
		}

		i++
		if i%10000 == 0 {
			log.Printf("Processed %d records", i)
		}
	}

	return ribEntries, nil
}
