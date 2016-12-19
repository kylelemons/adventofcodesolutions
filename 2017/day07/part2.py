#!/usr/bin/python2.7

import re
from itertools import chain

def isABA(s):
    return len(s) == 3 and s[0] == s[2] and s[0] != s[1]

def findABAs(s):
    abas = []
    for i in xrange(len(s)-2):
        ss = s[i:][:3]
        if isABA(ss):
            abas.append(ss)
    return abas

def bab(aba):
    return aba[1] + aba[0] + aba[1]

def supportsSSL(ip):
    groups = re.split(r'\[([a-z]*)\]', ip)
    abas = set(chain.from_iterable(findABAs(s) for s in groups[0::2]))
    babs = set([bab(aba) for aba in chain.from_iterable(findABAs(h) for h in groups[1::2])])
    return len(babs & abas) > 0

for ip, want in [('aba[bab]xyz', True),
                 ('xyx[xyx]xyx', False),
                 ('aaa[kek]eke', True),
                 ('zazbz[bzb]cdb', True)]:
    print ip, supportsSSL(ip), 'want', want

print sum([supportsSSL(ip) for ip in open('input.txt').readlines()])
