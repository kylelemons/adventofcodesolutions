#!/usr/bin/python2.7

from re import *
from itertools import *
from functools import *
from md5 import *
from math import *
from copy import *

def droptime(discs):
    discs, t = [(c,(p+i+1)%c) for i, (c, p) in enumerate(discs)], 0
    while True:
        print "AT t =", t, "SLOTS WILL BE", discs
        if all(p==0 for c, p in discs):
            return t
        discs, t = [(c,(p+1)%c) for c, p in discs], t + 1

tests = [
    # Disc #1 has 5 positions; at time=0, it is at position 4.
    # Disc #2 has 2 positions; at time=0, it is at position 1.
    (( [(5,4),(2,1)] ,), 5 ),
    (( [(5,4),(5,3),(5,2),(5,1),(5,0),(5,4)] ,), 0),
]
for args, ans in tests:
    got = droptime(*args)
    print args, got, ans, got == ans

# Disc #1 has 7 positions; at time=0, it is at position 0.
# Disc #2 has 13 positions; at time=0, it is at position 0.
# Disc #3 has 3 positions; at time=0, it is at position 2.
# Disc #4 has 5 positions; at time=0, it is at position 2.
# Disc #5 has 17 positions; at time=0, it is at position 0.
# Disc #6 has 19 positions; at time=0, it is at position 7.
print droptime([(7,0),(13,0),(3,2),(5,2),(17,0),(19,7)])
