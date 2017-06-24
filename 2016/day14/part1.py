#!/usr/bin/python2.7

import hashlib
import itertools

def memoize(f):
    cache = {}
    def wrapper(*args):
        key = tuple(args)
        if key in cache:
            return cache[key]
        res = f(*args)
        cache[key] = res
        return res
    return wrapper

@memoize
def candidate(salt, index):
    return hashlib.md5(salt+str(index)).hexdigest()

def has_repeat(h, length):
    for k, g in itertools.groupby(h):
        if len(list(g)) >= length:
            return k
    return None

def is_key(c, index):
    r = has_repeat(c, 3)
    if r != None:
        print "triple", index, c, r
        penta = r*5
        for i in xrange(1000):
            cc = candidate(salt, index+1+i)
            if penta in cc:
                print "penta ", index, cc, penta
                return True
    return False

index = 0
salt = 'yjdafjpo'
keys = []
while len(keys) < 64:
    c = candidate(salt, index)
    if is_key(c, index):
        keys.append((index, c))
    index += 1
