#!/usr/bin/python2.7

import re
import itertools
import functools
import md5
import math

df = open('input.txt').readlines()[2:]
df = [line.split() for line in df]

File = 0
Size = 1
Used = 2
Avail = 3
UseP = 4

def size(s):
    return int(s.strip('T'))

count = 0
for a in df:
    if a[Used] == "0T":
        continue
    for b in df:
        if a[File] == b[File]:
            continue
        if size(a[Used]) <= size(b[Avail]):
            count += 1
print count
