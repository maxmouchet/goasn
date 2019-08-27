package goasn

import (
	"compress/bzip2"
	"compress/gzip"
	"io"
	"net"
	"os"
	"strings"
)

// Similar to https://github.com/kaorimatz/go-mrt/blob/master/bgp.go#L29-L32
const (
	ASPathSegmentTypeSet = iota
	ASPathSegmentTypeSequence
)

type ASPathSegment struct {
	Type uint8
	ASNs []uint32
}

type RIBEntry struct {
	OriginCode string
	Network    net.IPNet
	LocPrf     uint32
	Path       []ASPathSegment
}

func RIBFromMRT(path string) ([]RIBEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var r io.Reader = file

	if strings.HasSuffix(path, ".bz2") {
		r = bzip2.NewReader(file)
	}

	if strings.HasSuffix(path, ".gz") {
		r, err = gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
	}

	return ribEntriesFromMRT(r)
}
