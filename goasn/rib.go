package goasn

import (
	"compress/bzip2"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
)

// Public API

type RIBEntry struct {
	OriginCode string
	Network    net.IPNet
	LocPrf     uint32
	Path       []uint32
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

	return ribEntriesFromMRT(r)
}

func RIBFromShowIPDump(path string) ([]RIBEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var r io.Reader = file
	if strings.HasSuffix(path, ".bz2") {
		r = bzip2.NewReader(file)
	}

	return ribEntriesFromShowIPDump(r)
}

// GetShowIPDumpRUL from Route Views for time `t`
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
