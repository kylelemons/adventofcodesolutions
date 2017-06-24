#!/usr/bin/python2.7

# NOTE: this doesn't actually do anything, I figured it out by hand...

import re
import itertools
import functools
import md5
import math
from heapq import *
import copy

df = open('input.txt').readlines()[2:]
df = [line.split() for line in df]

File = 0
Size = 1
Used = 2
Avail = 3
UseP = 4

def size(s):
    return int(s.strip('T'))

Width = 37
Height = 25
grid = [[df[x*25+y] for x in xrange(37)] for y in xrange(25)]

print "min size", min([min([size(n[Size]) for n in row]) for row in grid])
for row in grid:
    for col in row:
        if size(col[Used]) > 85:
            print col

# every pos or delta is (y,x)
deltas = [(-1, 0), (0, -1), (1, 0), (0, 1)]

# q contains (distance of data to dest, number of moves, location, state)
qDist = 0
qMoves = 1
qLoc = 2
qState = 3

if False:
    q = [(Width-1, 0, (0, Width-1), grid)]
    visited = {}
    while len(q) > 0:
        curr = heappop(q)
        loc = curr[qLoc]
        if loc[0] < 0 or loc[1] < 0:
            continue
        if loc[0] >= Height or loc[1] >= Width:
            continue
        if loc in visited:
            continue
        visited[loc] = True

        for i, dd in deltas:
            nextPos = curr[qLoc]
            nextPos[0] += dd[0]
            nextPos[1] += dd[1]
            ngrid = copy.deepcopy(curr[qState])
            heappush(q, (nextPos[0] + nextPos[1], curr[qMoves]+1, nextPos, ngrid))
