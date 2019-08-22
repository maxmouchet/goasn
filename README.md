# goasn

[![Build Status](https://travis-ci.org/maxmouchet/goasn.svg?branch=master)](https://travis-ci.org/maxmouchet/goasn)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxmouchet/goasn)](https://goreportcard.com/report/github.com/maxmouchet/goasn)

goasn provides fast lookup of IP addresses to AS numbers from BGP archives.  
It reads [pyasn](https://github.com/hadiasghari/pyasn) data files and store IP addresses in a radix tree for fast lookups.

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

ip := net.ParseIP("8.8.8.8")
asn, _ := asndb.LookupIP(ip)
// => 13335
```
