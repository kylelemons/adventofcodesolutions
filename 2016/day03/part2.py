#!/usr/bin/python2.7
from itertools import izip_longest


def tri(s):
    sides = sorted([int(s) for s in s])
    return 1 if sides[0] + sides[1] > sides[2] else 0

def tris(s):
    lines = iter(s.splitlines())
    count = 0
    try:
        while True:
            t0 = lines.next().strip().split()
            t1 = lines.next().strip().split()
            t2 = lines.next().strip().split()
            count += sum([tri([t0[i], t1[i], t2[i]]) for i in xrange(3)])
    except StopIteration:
        return count
    raise 'wtf?'

def do(s):
    print 'tris(' + s + ')'
    print ' =', tris(s)

do('5 10 25\n1 2 100\n100 10 99')
do(open('input1.txt').read())
