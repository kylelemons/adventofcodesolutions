#!/usr/bin/python2.7

import re
import itertools
import functools
import md5
import math

def lcr(prev, index):
    if index == 0:
        return False, prev[0], prev[1]
    if index == len(prev)-1:
        return prev[-2], prev[-1], False
    return prev[index-1], prev[index], prev[index+1]

def traps(row):
    return [c == '^' for c in row]

def chars(traps):
    return ''.join('^' if v else '.' for v in traps)

def is_trap(prev, index):
    left, center, right = lcr(prev, index)
    if left and center and not right:
        return True
    if center and right and not left:
        return True
    if left and not center and not right:
        return True
    if right and not center and not left:
        return True
    return False

def rows(start, count):
    last = traps(start)
    for i in xrange(count):
        yield last
        last = [is_trap(last, index) for index in xrange(len(last))]

def count(start, count):
    return sum(sum(not x for x in row) for row in rows(start, count))

tests = [
    (('..^^.', 3), 6),
    (('.^^.^.^^^^', 10), 38),
]
for args, ans in tests:
    got = count(*args)
    print args, got, ans, got == ans

print count('.^..^....^....^^.^^.^.^^.^.....^.^..^...^^^^^^.^^^^.^.^^^^^^^.^^^^^..^.^^^.^^..^.^^.^....^.^...^^.^.', 400000)
