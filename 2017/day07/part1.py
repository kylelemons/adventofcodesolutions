#!/usr/bin/python2.7

import re

def hasABBA(s):
    if len(s) < 4:
        return False
    if len(s) == 4:
        # also amusing, s[-1::-1] is the reversed s
        return s[0] == s[3] and s[1] == s[2] and s[0] != s[1]
    for i in xrange(len(s)-4):
        if hasABBA(s[i:i+4]):
            return True
    return False

for i, o in [('abba', True), ('abcd', False), ('ioxxoj', True)]:
    print i, hasABBA(i), 'want', o

def supportsTLS(ip):
    groups = re.split(r'\[([a-z]*)\]', ip)
    segments = [hasABBA(s) for s in groups[0::2]]
    hypernet = [hasABBA(h) for h in groups[1::2]]
    return any(segments) and not any(hypernet)

for ip, want in [('abba[mnop]qrst', True),
           ('abcd[bddb]xyyx', False),
           ('aaaa[qwer]tyui', False),
           ('ioxxoj[asdfgh]zxcvbn', True),
           ('[asdfgh]zxcvbn', False)]:
    print ip, supportsTLS(ip), 'want', want

print sum([supportsTLS(ip) for ip in open('input.txt').readlines()])
