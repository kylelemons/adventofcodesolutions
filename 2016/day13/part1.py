#!/usr/bin/python2.7

from re import *
from itertools import *
from functools import *
from md5 import *
from math import *
from copy import *

def pathlen(fav, dst):
    src = (1,1)

    def open(x,y):
        v = x*x + 3*x + 2*x*y + y + y*y + fav
        return bin(v).count('1')%2 == 0

    def ch(x,y, show):
        if not open(x,y):
            return '#'
        if (x,y) == src:
            return 'S'
        if (x,y) == dst:
            return 'X'
        if (x,y) in show:
            return 'O'
        return '.'

    def disp(show=set()):
        for y in xrange(dst[1]*3/2):
            print ''.join(ch(x,y,show) for x in xrange(dst[0]*3/2))
        print

    DELTAS = [
        (-1, 0), # left
        (+1, 0), # right
        (0, -1), # up
        (0, +1), # down
    ]
    MAX_DEPTH = 100

    disp()
    q = []
    #        ( path  )
    q.append(( [src] ))
    visited = set()

    while len(q) > 0:
        path = q.pop(0)
        loc = path[-1]
        if loc in visited:
            continue
        visited.add(loc)

        if loc == dst:
            disp(set(path))
            return len(path)-1

        if len(path) > MAX_DEPTH:
            continue

        for dd in DELTAS:
            nloc = tuple(sum(x) for x in izip(loc, dd))
            npath = path + [nloc]
            if not open(*nloc):
                continue
            q.append(( npath ))

    return None

tests = [
    ((10, (7,4)), 11)
]
for args, ans in tests:
    got = pathlen(*args)
    print args, got, ans, got == ans

print pathlen(1364, (31,39))
