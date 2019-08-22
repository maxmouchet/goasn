package main

import (
	"fmt"
	"net"
	"time"

	"github.com/maxmouchet/goasn"
)

var dbPath = "ipasn_20190822.dat"

func main() {
	fmt.Printf("Loading database from %s...\n", dbPath)
	asndb, err := goasn.NewDB(dbPath)

	start := time.Now()
	asn, err := asndb.LookupStr("8.8.8.8")
	fmt.Printf("Lookup took %s\n", time.Since(start))
	fmt.Println(asn, err)

	start = time.Now()
	ip := net.ParseIP("1.1.1.1")
	asn, err = asndb.LookupIP(ip)
	fmt.Printf("Lookup took %s\n", time.Since(start))
	fmt.Println(asn, err)
}
