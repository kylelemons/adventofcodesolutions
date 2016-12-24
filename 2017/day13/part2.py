#!/usr/bin/python2.7

from re import *
from itertools import *
from functools import *
from md5 import *
from math import *
from copy import *

def visitable(fav, MAX_DEPTH):
    src = (1,1)

    def open(x,y):
        v = x*x + 3*x + 2*x*y + y + y*y + fav
        return bin(v).count('1')%2 == 0

    def ch(x,y, show):
        if not open(x,y):
            return '#'
        if (x,y) == src:
            return 'S'
        if (x,y) in show:
            return 'O'
        return '.'

    def disp(show=set()):
        for y in xrange(MAX_DEPTH):
            print ''.join(ch(x,y,show) for x in xrange(MAX_DEPTH))
        print

    DELTAS = [
        (-1, 0), # left
        (+1, 0), # right
        (0, -1), # up
        (0, +1), # down
    ]

    q = []
    #        ( path  )
    q.append(( [src] ))
    visited = set()

    while len(q) > 0:
        path = q.pop(0)
        loc = path[-1]
        if loc[0] < 0 or loc[1] < 0:
            continue
        if loc in visited:
            continue
        visited.add(loc)

        if len(path)-1 >= MAX_DEPTH:
            continue

        for dd in DELTAS:
            nloc = tuple(sum(x) for x in izip(loc, dd))
            npath = path + [nloc]
            if not open(*nloc):
                continue
            q.append(( npath ))

    disp(visited)
    return len(visited)

print visitable(1364, 12)
print visitable(1364, 50)
