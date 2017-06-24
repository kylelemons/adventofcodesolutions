#!/usr/bin/python2.7

import md5

base = md5.new()
base.update('reyedfim')

def withcounter(c):
    m = base.copy()
    m.update(str(c))
    h = m.hexdigest()
    if h.startswith('00000'):
        return h[5]
    return None

code = []
counter = 0
while (len(code) < 8):
    nxt = withcounter(counter)
    if nxt:
        print nxt
        code += [nxt]
    counter+=1
print ''.join(code)
