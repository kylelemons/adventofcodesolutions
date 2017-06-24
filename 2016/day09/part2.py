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
        chars, count = int(m.group(1)), int(m.group(2))
        marker, marked = m.group(0), cur[len(m.group(0)):][:chars]
        i, length = i+len(marker)+chars, length+decompress(marked)*count
    return length

print decompress('(27x12)(20x12)(13x14)(7x10)(1x12)A')
print decompress(open('input.txt').read().strip())
