#!/usr/bin/python2.7

import re
import itertools
import functools
import md5
import math

blacklist_text = """5-8
0-2
4-7"""

blacklist_text = open('input.txt').read()

blacklist = [[int(x) for x in line.split('-')] for line in blacklist_text.splitlines()]
blacklist = sorted(blacklist)

def combined(blacklist):
    lo, hi = blacklist.pop(0)
    for clo, chi in blacklist:
        if clo > hi:
            yield (lo, hi)
            lo, hi = clo, chi
            continue
        if chi > hi:
            hi = chi
    yield (lo, hi)

nextPassed = 0
for lo, hi in combined(blacklist):
    print lo, hi
    if nextPassed < lo:
        print nextPassed
        break
    nextPassed = hi+1
