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
	ribEntry.Path = make([]ASPathSegment, 0)

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
				segment := ASPathSegment{}
				segment.ASNs = make([]uint32, len(asPathSegment.Value))

				switch asPathSegment.Type {
				case mrt.BGPASPathSegmentTypeASSet:
					segment.Type = ASPathSegmentTypeSet
				case mrt.BGPASPathSegmentTypeASSequence:
					segment.Type = ASPathSegmentTypeSequence
				}

				for i, v := range asPathSegment.Value {
					switch len(v) {
					case 2:
						segment.ASNs[i] = uint32(binary.BigEndian.Uint16(v))
					case 4:
						segment.ASNs[i] = binary.BigEndian.Uint32(v)
					}
				}

				ribEntry.Path = append(ribEntry.Path, segment)
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
			log.Printf("Error: %s", err)
			continue
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
