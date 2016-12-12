#!/usr/bin/python2.7

import re
import itertools

def checksum(name):
    counts = [(k, len(list(l))) for k, l in itertools.groupby(sorted(name)) if k != '-']
    return ''.join([p[0] for p in sorted(counts, reverse=True, key=lambda p: p[1])[:5]])

def parse(room):
    return re.match(r"(?P<name>[-a-z]+)-(?P<sector>[0-9]+)\[(?P<checksum>[-a-z]+)\]", room).groupdict()

def valid(r):
    return checksum(r['name']) == r['checksum']

a = ord('a')
def rotate(c, by):
    return chr(a+(ord(c)-a+by)%26) if c != '-' else ' '

def decrypt(r):
    by = int(r['sector'])
    return ''.join([rotate(c, by) for c in r['name']])

def decryptall(rooms):
    return {decrypt(r): int(r['sector']) for r in [parse(room) for room in rooms] if valid(r)}

def do(rooms):
    rooms = [room.strip() for room in rooms]
    for room, sector in sorted(decryptall(rooms).iteritems()):
        if re.match(r"north.*pole", room): # comment for your own amusement
            print room, sector

do(['aaaaa-bbb-z-y-x-123[abxyz]'])
do(['a-b-c-d-e-f-g-h-987[abcde]'])
do(['not-a-real-room-404[oarel]'])
do(['totally-real-room-200[decoy]'])
do(['qzmt-zixmtkozy-ivhz-343['+checksum('qzmt-zixmtkozy-ivhz')+']'])
do(open('input.txt').readlines())
