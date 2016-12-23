#!/usr/bin/python2.7

from re import *
from itertools import *
from functools import *
from md5 import *
from math import *
from copy import *

def answer():
    pass

tests = [
    ((), None)
]
for args, ans in tests:
    got = answer(*args)
    print args, got, ans, got == ans
