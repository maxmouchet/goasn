package goasn

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
)

type RIBEntry struct {
	OriginCode string
	StatusCode string
	Network    net.IPNet
	NextHop    net.IP
	Metric     uint32
	LocPrf     uint32
	Weight     uint32
	Path       []uint32
}

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

	metric, err := strconv.ParseUint(matches[4], 10, 32)
	if err != nil {
		return nil, err
	}

	locPrf, err := strconv.ParseUint(matches[5], 10, 32)
	if err != nil {
		return nil, err
	}

	weight, err := strconv.ParseUint(matches[6], 10, 32)
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
		StatusCode: matches[1],
		Network:    *network,
		NextHop:    net.ParseIP(matches[3]),
		Metric:     uint32(metric),
		LocPrf:     uint32(locPrf),
		Weight:     uint32(weight),
		Path:       path,
	}, nil
}

func ReadShowIPDump(path string) ([]RIBEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var r io.Reader = file
	if strings.HasSuffix(path, ".bz2") {
		r = bzip2.NewReader(file)
	}

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

// TODO: MRT Format

// http://archive.routeviews.org/oix-route-views/2019.08/oix-full-snapshot-2019-08-02-0600.bz2
// TODO Old show ip bgp format (.dat files)

func GetShowIPDumpURL(t time.Time) string {
	return fmt.Sprintf(
		"%s/%s/oix-full-snapshot-%s.bz2",
		"http://archive.routeviews.org/oix-route-views",
		t.UTC().Format("2006.01"),
		t.UTC().Format("2006-01-02-1500"),
	)
}

// DownloadShowIPDump from Route Views for time `t` to directory `dir`
func DownloadShowIPDump(dir string, t time.Time) (string, error) {
	url := GetShowIPDumpURL(t)
	dst := filepath.Join(dir, filepath.Base(url))

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	log.Printf("Downloading %s to %s...", url, dst)
	_, err = grab.Get(dst, url)

	return dst, err
}
