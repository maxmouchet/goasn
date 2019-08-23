package goasn

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// `sh ip bgp format`
// TODO: Old show ip bgp format (.dat files)

var showIPLinePattern = regexp.MustCompile(`^([dhirsS*>])\s+([0-9./]+)\s+([0-9.]+)\s+(\d+)\s+(\d+)\s+(\d+)\s+([\d\s]+)\s+([ie?])$`)

func parseShowIPLine(line string) (*RIBEntry, error) {
	matches := showIPLinePattern.FindStringSubmatch(line)
	if len(matches) != 9 {
		return nil, fmt.Errorf("Failed to parse line %s", line)
	}

	_, network, err := net.ParseCIDR(matches[2])
	if err != nil {
		return nil, err
	}

	locPrf, err := strconv.ParseUint(matches[5], 10, 32)
	if err != nil {
		return nil, err
	}

	pathStr := strings.Split(matches[7], " ")
	path := make([]uint32, len(pathStr))
	for i, asnStr := range pathStr {
		asn, err := strconv.ParseUint(asnStr, 10, 32)
		if err != nil {
			return nil, err
		}
		path[i] = uint32(asn)
	}

	return &RIBEntry{
		OriginCode: matches[8],
		Network:    *network,
		LocPrf:     uint32(locPrf),
		Path:       path,
	}, nil
}

func ribEntriesFromShowIPDump(r io.Reader) ([]RIBEntry, error) {
	entries := make([]RIBEntry, 0)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ";") {
			continue
		}

		entry, err := parseShowIPLine(line)
		if err != nil {
			log.Println(err)
		} else {
			entries = append(entries, *entry)
		}
	}

	return entries, nil
}
