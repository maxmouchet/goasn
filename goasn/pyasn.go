package goasn

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type PyASNEntry struct {
	IPNet net.IPNet
	ASN   uint32
}

var linePattern = regexp.MustCompile(`^(.+?)\t(\d+)$`)

func parseLine(line string) (*PyASNEntry, error) {
	matches := linePattern.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, fmt.Errorf("Failed to parse line %s", line)
	}

	_, ipn, err := net.ParseCIDR(matches[1])
	if err != nil {
		return nil, err
	}

	asn, err := strconv.ParseUint(matches[2], 10, 32)
	if err != nil {
		return nil, err
	}

	entry := PyASNEntry{
		IPNet: *ipn,
		ASN:   uint32(asn),
	}

	return &entry, nil
}

func ReadPyASNDatabase(path string) ([]PyASNEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries := make([]PyASNEntry, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ";") {
			continue
		}

		entry, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		entries = append(entries, *entry)
	}

	return entries, nil
}
