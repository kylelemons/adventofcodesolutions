#!/usr/bin/python2.7

import md5

base = md5.new()
base.update('reyedfim')

def withcounter(c):
    m = base.copy()
    m.update(str(c))
    h = m.hexdigest()
    if h.startswith('00000'):
        return h[5], h[6]
    return None, None

code = ['_']*8
counter = 0
found = 0
while found < 8:
    pos, ch = withcounter(counter)
    if pos and pos >= '0' and pos <= '7':
        pos = int(pos)
        if code[pos] == '_':
            code[pos] = ch
            print ''.join(code)
            found += 1
    counter+=1
