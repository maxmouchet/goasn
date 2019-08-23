package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/maxmouchet/goasn/goasn"
)

var addrsPath = "addrs.txt"
var dbPath = "ipasn_20190822.dat"

func readAddrsAsStr(path string) []string {
	file, _ := os.Open(path)
	defer file.Close()

	addrs := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		addrs = append(addrs, scanner.Text())
	}

	return addrs
}

func readAddrsAsIP(path string) []net.IP {
	addrsStr := readAddrsAsStr(path)
	addrsIP := make([]net.IP, len(addrsStr))
	for i, addr := range addrsStr {
		addrsIP[i] = net.ParseIP(addr)
	}
	return addrsIP
}

func main() {
	fmt.Printf("Loading database from %s...\n", dbPath)
	start := time.Now()
	asndb, _ := goasn.NewDB(dbPath)
	fmt.Printf("Took %s\n\n", time.Since(start))

	addrsStr := readAddrsAsStr(addrsPath)
	addrsIP := readAddrsAsIP(addrsPath)

	fmt.Printf("Looking up %d addresses (string)\n", len(addrsStr))
	start = time.Now()
	for _, addr := range addrsStr {
		asndb.LookupStr(addr)
	}
	fmt.Printf("Took %s\n\n", time.Since(start))

	fmt.Printf("Looking up %d addresses (net.IP)\n", len(addrsIP))
	start = time.Now()
	for _, addr := range addrsIP {
		asndb.LookupIP(addr)
	}
	fmt.Printf("Took %s\n", time.Since(start))
}
