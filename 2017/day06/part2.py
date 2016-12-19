#!/usr/bin/python2.7

from itertools import groupby

def decode(lines):
    cols = zip(*lines)
    scols = [sorted(col) for col in cols]
    freqs = [[(len(list(v)), k) for k, v in groupby(col)] for col in scols]
    letters = [sorted(col)[0][1] for col in freqs]
    return ''.join(letters)

print decode("""eedadn
drvtee
eandsr
raavrd
atevrs
tsrnev
sdttsa
rasrtv
nssdts
ntnada
svetve
tesnvt
vntsnd
vrdear
dvrsen
enarar""".splitlines())

print decode(open('input.txt').readlines())
