#!/usr/bin/python2.7

import re

def decompress(s):
    i, length = 0, 0
    while i < len(s):
        cur = s[i:]
        m = re.match(r'\((\d+)x(\d+)\)', cur)
        if not m:
            i, length = i+1, length+1
            continue
        marker, chars, count = m.group(0), int(m.group(1)), int(m.group(2))
        i, length = i+len(marker)+chars, length+chars*count
    return length

print decompress('ADVENT'), 6
print decompress('A(1x5)BC'), 7
print decompress('(3x3)XYZ'), 9
print decompress('A(2x2)BCD(2x2)EFG'), 11
print decompress('(6x1)(1x3)A'), 6
print decompress('X(8x2)(3x3)ABCY'), 18
print decompress(open('input.txt').read().strip())
