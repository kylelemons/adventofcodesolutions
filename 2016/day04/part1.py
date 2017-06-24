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

def sectorsum(rooms):
    return sum([int(r['sector']) for r in [parse(room) for room in rooms] if valid(r)])

def do(rooms):
    rooms = [room.strip() for room in rooms]
    print rooms
    print ' =', sectorsum(rooms)

do(['aaaaa-bbb-z-y-x-123[abxyz]'])
do(['a-b-c-d-e-f-g-h-987[abcde]'])
do(['not-a-real-room-404[oarel]'])
do(['totally-real-room-200[decoy]'])
do(open('input.txt').readlines())
