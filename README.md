# goasn

[![Build Status](https://travis-ci.org/maxmouchet/goasn.svg?branch=master)](https://travis-ci.org/maxmouchet/goasn)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxmouchet/goasn)](https://goreportcard.com/report/github.com/maxmouchet/goasn)

goasn provides fast lookup of IP addresses to AS numbers from BGP archives.  
It supports the following formats:
- [x] [pyasn](https://github.com/hadiasghari/pyasn) data files
- [ ] [Route Views](http://archive.routeviews.org/) `sh ip bgp` format RIBs
- [ ] [Route Views](http://archive.routeviews.org/) MRT format RIBs

It reads  and store IP addresses in a radix tree ([kentik/patricia](https://github.com/kentik/patricia)) for fast lookups.

## Quick Start

```bash
go get github.com/maxmouchet/goasn
```

From [pyasn](https://github.com/hadiasghari/pyasn) documentation:
```bash
pyasn_util_download.py --latest
pyasn_util_convert.py --single <Downloaded RIB File> <ipasn_db_file_name>
```

```go
asndb, _ := goasn.NewDB("ipasn_db.dat")

asn, _ := asndb.LookupStr("8.8.8.8")
// => 15169

ip := net.ParseIP("1.1.1.1")
asn, _ := asndb.LookupIP(ip)
// => 13335
```

## Performance

Constructing the radix tree is slower than pyasn, but lookups are faster.

```bash
# benchmark.go
Loading database from ipasn_20190822.dat...
Took 1.719656297s

Looking up 10000 addresses (string)
Took 3.367572ms

Looking up 10000 addresses (net.IP)
Took 1.631176ms
```

```bash
# benchmark.py
Loading database from ipasn_20190822.dat
Took 197.402488ms

Looking up 10000 addresses (string)
Took 9.340981ms
```
