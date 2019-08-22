#!/usr/bin/env python3

import pyasn
import time

addrs_path = "addrs.txt"
db_path = "ipasn_20190822.dat"

def read_addrs(path):
    return list(map(str.strip, open(path)))

if __name__ == '__main__':
    print('Loading database from {}'.format(db_path))
    start = time.time_ns()
    asndb = pyasn.pyasn(db_path)
    print('Took {}ms\n'.format((time.time_ns()-start)/1e6))

    addrs = read_addrs(addrs_path)

    print('Looking up {} addresses (string)'.format(len(addrs)))
    start = time.time_ns()
    for addr in addrs:
        asndb.lookup(addr)
    print('Took {}ms'.format((time.time_ns()-start)/1e6))