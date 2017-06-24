#!/usr/bin/python2.7

from re import *
from itertools import *
from functools import *
from md5 import *
from math import *
from copy import *
import string

def grouper(iterable, n, fillvalue=None):
    "Collect data into fixed-length chunks or blocks"
    # grouper('ABCDEFG', 3, 'x') --> ABC DEF Gxx
    args = [iter(iterable)] * n
    return izip_longest(fillvalue=fillvalue, *args)

def checksum(s):
    "Return the pairwise checksum of s"
    return ''.join('1' if a == b else '0' for a, b in grouper(s, 2, '0'))

def fill_checksum(base, length):
    "Fill the base data until it's the given length then checksum the data"
    invert = string.maketrans('01', '10')
    while len(base) < length:
        base = base + '0' + ''.join(reversed(base.translate(invert)))
    base = base[:length]
    chk = checksum(base)
    while len(chk) % 2 == 0:
        chk = checksum(chk)
    return chk

tests = [
    (('110010110100', 12), '100'),
    (('10000', 20), '01100'),
]
for args, ans in tests:
    got = fill_checksum(*args)
    print args, got, ans, got == ans

print fill_checksum('10001110011110000', 272)
