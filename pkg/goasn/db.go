package goasn

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// pyasn-style database

type ASNDatabase struct {
	Entries []PrefixOrigin
}

type IXPDatabase struct {
	Entries []PrefixIXP
}

var linePattern = regexp.MustCompile(`^(.+?)\t(.+)$`)

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

// TODO: Deduplicate code

func (p PrefixOrigin) MarshalText(singleAS bool) ([]byte, error) {
	asns := p.Origin
	if singleAS {
		asns = asns[0:1]
	}

	str := fmt.Sprintf(
		"%s\t%s\n",
		p.Prefix.String(),
		formatSlice(asns),
	)

	return []byte(str), nil
}

func (p PrefixIXP) MarshalText() ([]byte, error) {
	str := fmt.Sprintf(
		"%s\t%s\n",
		p.Prefix.String(),
		p.IXP,
	)
	return []byte(str), nil
}

func (p *PrefixOrigin) UnmarshalText(data []byte) error {
	str := string(data)
	matches := linePattern.FindStringSubmatch(str)
	if len(matches) != 3 {
		return fmt.Errorf("Failed to parse string %s", str)
	}

	_, prefix, err := net.ParseCIDR(matches[1])
	if err != nil {
		return err
	}

	asnsStr := strings.Split(matches[2], ",")
	asns := make([]uint32, len(asnsStr))
	for i, str := range asnsStr {
		asn, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return err
		}
		asns[i] = uint32(asn)
	}

	p.Prefix = *prefix
	p.Origin = asns

	return nil
}

func (p *PrefixIXP) UnmarshalText(data []byte) error {
	str := string(data)
	matches := linePattern.FindStringSubmatch(str)
	if len(matches) != 3 {
		return fmt.Errorf("Failed to parse string %s", str)
	}

	_, prefix, err := net.ParseCIDR(matches[1])
	if err != nil {
		return err
	}

	p.Prefix = *prefix
	p.IXP = matches[2]

	return nil
}

// TODO: Deduplicate code
func (db ASNDatabase) MarshalText(singleAS bool) ([]byte, error) {
	w := new(bytes.Buffer)

	// TODO: Diff with pyasn
	// TODO: Original source, Prefixes-v4/v6
	fmt.Fprintf(w, "; IP-ASN32-DAT file\n")
	fmt.Fprintf(w, "; Original source:%s\n")
	fmt.Fprintf(w, "; Converted on:\t%s\n", time.Now().Format("Mon Jan 2 15:04:05 2006"))
	fmt.Fprintf(w, "; Prefixes-v4:\t%d\n")
	fmt.Fprintf(w, "; Prefixes-v6:\t%d\n")
	fmt.Fprintf(w, ";\n")

	lastNet := ""

	_, defaultV4, _ := net.ParseCIDR("0.0.0.0/0")
	_, defaultV6, _ := net.ParseCIDR("::/0")

	// WARN if same prefix with differents ASes

	for _, entry := range db.Entries {
		if entry.Prefix.String() == lastNet {
			continue
		}

		// TODO: Optimize
		if (entry.Prefix.String() == defaultV4.String()) || (entry.Prefix.String() == defaultV6.String()) {
			continue
		}

		lastNet = entry.Prefix.String()

		b, err := entry.MarshalText(singleAS)
		if err != nil {
			return nil, err
		}

		_, err = w.Write(b)
		if err != nil {
			return nil, err
		}
	}

	return w.Bytes(), nil
}

func (db *ASNDatabase) UnmarshalText(data []byte) error {
	// TODO: Cleanup/optimize
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	entries := make([]PrefixOrigin, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ";") {
			continue
		}

		var p PrefixOrigin
		err := p.UnmarshalText([]byte(line))
		if err != nil {
			return err
		}

		entries = append(entries, p)
	}

	db.Entries = entries
	return nil
}

// TODO: Deduplicate code
func (db IXPDatabase) MarshalText() ([]byte, error) {
	w := new(bytes.Buffer)

	// TODO: Diff with pyasn
	// TODO: Original source, Prefixes-v4/v6
	fmt.Fprintf(w, "; IP-IXP-DAT file\n")
	fmt.Fprintf(w, "; Original source:%s\n")
	fmt.Fprintf(w, "; Converted on:\t%s\n", time.Now().Format("Mon Jan 2 15:04:05 2006"))
	fmt.Fprintf(w, "; Prefixes-v4:\t%d\n")
	fmt.Fprintf(w, "; Prefixes-v6:\t%d\n")
	fmt.Fprintf(w, ";\n")

	for _, entry := range db.Entries {
		b, err := entry.MarshalText()
		if err != nil {
			return nil, err
		}
		_, err = w.Write(b)
		if err != nil {
			return nil, err
		}
	}

	return w.Bytes(), nil
}

func (db *IXPDatabase) UnmarshalText(data []byte) error {
	// TODO: Cleanup/optimize
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	entries := make([]PrefixIXP, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ";") {
			continue
		}

		var p PrefixIXP
		err := p.UnmarshalText([]byte(line))
		if err != nil {
			return err
		}

		entries = append(entries, p)
	}

	db.Entries = entries
	return nil
}
