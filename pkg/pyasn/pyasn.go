package pyasn

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
	ASNs  []uint32
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
		ASNs:  []uint32{uint32(asn)},
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

func WritePyASNDatabase(path string, entries []PyASNEntry) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	w.WriteString("; IP-ASN32-DAT file\n")
	w.WriteString("; Extended version with multiple ASes\n")
	w.WriteString(";\n")

	lastNet := ""

	// WARN if same prefix with differents ASes

	for _, entry := range entries {
		if entry.IPNet.String() == lastNet {
			continue
		}
		lastNet = entry.IPNet.String()
		fmt.Fprintf(
			w,
			"%s\t%d\n",
			entry.IPNet.String(),
			entry.ASNs,
		)
	}

	return nil
}
