#!/usr/bin/python2.7

import hashlib

Width = 4
Height = 4

Directions = {
    'U': (0, -1),
    'D': (0, +1),
    'L': (-1, 0),
    'R': (+1, 0),
}

def doors(password, path):
    path = ''.join(path)
    hexes = hashlib.md5(password+path).hexdigest()
    return {
        'U': hexes[0] >= 'b',
        'D': hexes[1] >= 'b',
        'L': hexes[2] >= 'b',
        'R': hexes[3] >= 'b',
    }

def path(password):
    # q contains (dist, pos[x,y], path)
    q = []
    q.append(([0,0], []))
    while len(q) > 0:
        pos, path = q.pop(0)
        if pos[0] < 0 or pos[1] < 0 or pos[0] >= Width or pos[1] >= Height:
            continue
        #print pos, path
        if pos == [3,3]:
            return ''.join(path)
            break
        if len(path) > 100:
            break
        opens = doors(password, path)
        for ch, dd in Directions.iteritems():
            if not opens[ch]:
                continue
            npos = [pos[i] + delta for i, delta in enumerate(dd)]
            npath = path + [ch]
            q.append((npos, npath))
    return 'AAAAAH'

tests = [
    ('ihgpwlah', 'DDRRRD'),
    ('kglvqrro', 'DDUDRLRRUDRD'),
    ('ulqzkmiv', 'DRURDRUDDLLDLUURRDULRLDUUDDDRR'),
]
for args, ans in tests:
    got = path(args)
    print args, got, ans, got == ans

print
print path('vkjiggvb')
