#!/usr/bin/python2.7

def tri(s):
    sides = sorted([int(s) for s in s.strip().split()])
    print sides, sides[0] + sides[1], 'valid' if sides[0] + sides[1] > sides[2] else 'invalid'
    return 1 if sides[0] + sides[1] > sides[2] else 0

def tris(s):
    return sum([tri(s) for s in s.splitlines()])

def do(s):
    print 'tris(' + s + ')'
    print ' =', tris(s)

do('5 10 25')
do(open('input1.txt').read())
