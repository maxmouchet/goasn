# goasn

[![Build Status](https://travis-ci.org/maxmouchet/goasn.svg?branch=master)](https://travis-ci.org/maxmouchet/goasn)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxmouchet/goasn)](https://goreportcard.com/report/github.com/maxmouchet/goasn)
[![GoDoc](https://godoc.org/github.com/maxmouchet/goasn?status.svg)](https://godoc.org/github.com/maxmouchet/goasn) 

goasn provides fast lookup of IP addresses to IXP and AS numbers from RIB archives.  
It supports the following  sources:
- [Route Views](http://archive.routeviews.org/)/[RIS](https://www.ripe.net/analyse/internet-measurements/routing-information-service-ris/routing-information-service-ris) MRT format RIBs
- [PeeringDB](https://peeringdb.com/) data
- [pyasn](https://github.com/hadiasghari/pyasn) data files

It reads  and store IP addresses in a radix tree ([kentik/patricia](https://github.com/kentik/patricia)) for fast lookups.

## Quick Start


### CLI

You can download the [latest binary](https://github.com/maxmouchet/goasn/releases) or build goasn by yourself.  
Building meshmon requires [Go](https://golang.org/dl/) 1.12+.
```bash
git clone https://github.com/maxmouchet/goasn.git
cd goasn; make
```

```bash
goasn download --collector route-views.amsix.routeviews.org --date 2019-08-01T08:00
goasn convert rib.20190801.0800.bz2
goasn lookup --db rib.20180801.0800.txt 8.8.8.8
```

### Library

```bash
go get github.com/maxmouchet/goasn
```

```go
asndb, _ := goasn.NewDB("rib.20180801.0800.txt")

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
