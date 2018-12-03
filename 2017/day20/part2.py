import re
from math import copysign

input = open('input.txt').readlines()

# list of (manhat, pos, vel, accel)
particles = []

linere = re.compile(r'p=<(-?\d+),(-?\d+),(-?\d+)>, v=<(-?\d+),(-?\d+),(-?\d+)>, a=<(-?\d+),(-?\d+),(-?\d+)>')
for line in input:
    m = linere.match(line)
    if not m:
        continue
    g = m.groups()
    p = [int(n) for n in g[0:][:3]]
    v = [int(n) for n in g[3:][:3]]
    a = [int(n) for n in g[6:][:3]]
    particles.append((p,v,a))

def signeq(p, v, a):
    if a == 0:
        a = copysign(0.0, v) # either +0 or -0
    return copysign(1, p) == copysign(1, v) == copysign(1, a)

def part1(parts):
    rounds = 100

    by_loc = {tuple(pp): (pp, vv, aa) for (pp, vv, aa) in parts}
    for round in range(rounds):
        seen = {}
        new_by_loc = {}
        for k, (pp, vv, aa) in by_loc.iteritems():
            vv = [sum(z) for z in zip(vv, aa)]
            pp = [sum(z) for z in zip(pp, vv)]

            if tuple(pp) not in seen:
                seen[tuple(pp)] = 0
            seen[tuple(pp)] += 1
            new_by_loc[tuple(pp)] = (pp, vv, aa)

        by_loc = {k: v for k, v in new_by_loc.iteritems() if seen[k] == 1}
        print round, len(by_loc)
    return len(by_loc)

print ('input', part1(particles))
